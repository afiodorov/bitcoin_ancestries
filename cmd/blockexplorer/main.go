package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"sort"

	"github.com/go-pg/pg"
)

type Graph map[string][]string
type Transaction struct {
	Hash        string
	BlockNumber int
	Children    []string
}
type Group struct {
	Transactions []Transaction
	AllChildren  []string
}

func processFile(file string, transactions chan<- Transaction) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalln(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(f)
	_, err = r.Read()
	if err != nil {
		log.Fatalln(err)
	}

	graph := make(Graph)
	ts := make([]string, 0, 1000)
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalln(err)
		}

		transaction := rec[0]
		input := rec[2]

		ts = append(ts, transaction)
		if val, ok := graph[transaction]; ok {
			graph[transaction] = append(val, input)
		} else {
			graph[transaction] = []string{input}
		}
	}
	f.Close()

	blockNumber := getBlockNumber(file)
	uniqueTransactions := uniqueOrdered(ts)
	for _, transaction := range uniqueTransactions {
		children, ok := graph[transaction]
		if !ok {
			log.Fatalf("Transaction %v not found in graph %v\n", string(transaction), graph)
		}
		transactions <- Transaction{
			Hash:        transaction,
			BlockNumber: blockNumber,
			Children:    children,
		}
	}

}

func groupTransactions(transactions <-chan Transaction, groups chan<- Group) {
	group := Group{}

	prevBlock := -1
	for t := range transactions {
		if prevBlock == -1 {
			prevBlock = t.BlockNumber
		}
		group.Transactions = append(group.Transactions, t)
		group.AllChildren = append(group.AllChildren, t.Children...)
		if t.BlockNumber-prevBlock >= fetchStep {
			prevBlock = t.BlockNumber
			group.AllChildren = unique(group.AllChildren)
			groups <- group
			group = Group{}
		}
	}
	group.AllChildren = unique(group.AllChildren)
	groups <- group

	close(groups)
}

func processTransactions(resHolder *ResHolder, db *pg.DB, groups <-chan Group, inserts chan<- map[string]Ancestry) {
	for g := range groups {
		resHolder.Fetch(db, g.AllChildren)
		var latestBlock int
		for _, t := range g.Transactions {
			out := makeAncestry(t.Hash, t.Children, t.BlockNumber, *resHolder)
			resHolder.m[t.Hash] = out
			latestBlock = t.BlockNumber
		}

		if len(resHolder.m) >= syncStep {
			log.Printf("Block %v. Inserting %v transactions...\n", latestBlock, len(resHolder.m))
			inserts <- resHolder.m
			resHolder.Clear()
		}
	}
	close(inserts)
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	var (
		pgHost     = flag.String("pghost", "localhost", "address of the postgress database")
		pgPort     = flag.Int("pgport", 5433, "port for ethereum postgres database")
		pgUser     = flag.String("pguser", "postgres", "user for ethereum postgres database")
		pgPassword = flag.String("pgpassword", "writer", "password for ethereum postgres database")
		pgDatabase = flag.String("pgdatabase", "postgres", "name of the database")
		sourceDir  = flag.String("s", "", "source directory")
		fromBlock  = flag.Int("f", 0, "resume from block")
	)
	flag.Parse()

	opts := &pg.Options{
		Addr:     fmt.Sprintf("%v:%d", *pgHost, *pgPort),
		User:     *pgUser,
		Password: *pgPassword,
		Database: *pgDatabase,
	}

	db := pg.Connect(opts)
	defer db.Close()

	err := createSchema(db)
	if err != nil {
		log.Fatal(err)
	}

	var files []string
	err = filepath.Walk(*sourceDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}

	sort.Slice(files, func(i, j int) bool {
		b1 := getBlockNumber(files[i])
		b2 := getBlockNumber(files[j])
		return b1 < b2
	})

	transactions := make(chan Transaction, 10000)
	go func() {
		for _, file := range files {
			if getBlockNumber(file) >= *fromBlock {
				processFile(file, transactions)
			}
		}
		close(transactions)
	}()

	groups := make(chan Group, 20)
	go groupTransactions(transactions, groups)

	inserts := make(chan map[string]Ancestry)
	rh := NewResHolder()
	go processTransactions(&rh, db, groups, inserts)

	for ins := range inserts {
		outs := make([]Ancestry, 0, len(ins))

		for _, v := range ins {
			outs = append(outs, v)
		}

		if _, err := db.Model(&outs).Insert(); err != nil {
			log.Fatalf("%v\n%v\n", err, outs)
		}
	}
}
