package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/afiodorov/bitcoin_ancestries/cmd/blocksorter/graph"
)

type Row []string

func (r Row) getEdge() (from, to string) {
	return r[2], r[0]
}

func chronSortRows(file string) []Row {
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

	rows := make([]Row, 0, 1000)
	rowNums := make(map[string][]int)
	g := graph.MakeGraph()

	rowNum := 0

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalln(err)
		}

		rows = append(rows, row)
		from, to := Row(row).getEdge()
		g.AddEdge(from, to)

		if nums, ok := rowNums[to]; ok {
			rowNums[to] = append(nums, rowNum)
		} else {
			rowNums[to] = []int{rowNum}
		}
		rowNum++
	}

	sorted, err := g.TSort()
	if err != nil {
		log.Fatalln(err)
	}

	res := make([]Row, 0, len(rows))

	for _, t := range sorted {
		for _, n := range rowNums[t] {
			res = append(res, rows[n])
		}
	}

	if len(res) != len(rows) {
		log.Fatalf("Missing rows. Old: %v, new: %v\n", len(rows), len(res))
	}

	f.Close()

	return res
}

func write(file string, rows []Row) {
	f, err := os.Create(file)
	if err != nil {
		log.Fatalln(err)
	}

	defer f.Close()
	w := csv.NewWriter(f)

	defer w.Flush()

	err = w.Write([]string{"hash", "value", "spent_transaction_hash"})
	if err != nil {
		log.Fatalln(err)
	}

	for _, row := range rows {
		err = w.Write(row)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func main() {
	sourceDir := flag.String("s", "", "source directory")
	destDir := flag.String("d", "", "destination directory")
	flag.Parse()

	err := os.MkdirAll(*destDir, os.ModePerm)

	if err != nil {
		log.Fatalln(err)
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

	maxGoroutines := 100
	guard := make(chan struct{}, maxGoroutines)

	var wg sync.WaitGroup

	for _, file := range files {
		guard <- struct{}{}

		wg.Add(1)

		go func(file string) {
			rows := chronSortRows(file)
			write(path.Join(*destDir, path.Base(file)), rows)
			log.Printf("Finished processing %v\n", file)
			<-guard
			wg.Done()
		}(file)
	}

	wg.Wait()
}
