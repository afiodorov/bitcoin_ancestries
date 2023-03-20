package main

import (
	"compress/gzip"
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
)

func write(file string, row []string) {
	var (
		f   *os.File
		err error
	)

	isNew := false

	if _, err := os.Stat(file); err == nil {
		f, err = os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		f, err = os.Create(file)
		if err != nil {
			log.Fatalln(err)
		}
		isNew = true
	}

	defer f.Close()
	w := csv.NewWriter(f)

	defer w.Flush()

	if isNew {
		err = w.Write([]string{"hash", "value", "spent_transaction_hash"})
		if err != nil {
			log.Fatalln(err)
		}
	}

	err = w.Write(row)
	if err != nil {
		log.Fatalln(err)
	}
}

func processFile(destDir, file string, wg *sync.WaitGroup, rowsGroup []chan []string) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalln(err)
	}

	gz, err := gzip.NewReader(f)

	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(gz)
	_, err = r.Read()

	if err != nil {
		log.Fatalln(err)
	}

	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalln(err)
		}

		blockNumber := rec[0]

		blockNumberI, err := strconv.Atoi(blockNumber)
		if err != nil {
			log.Fatalln(err)
		}

		row := []string{path.Join(destDir, blockNumber+".csv"), rec[1], rec[2], rec[3]}
		rowsGroup[blockNumberI%len(rowsGroup)] <- row
	}
	log.Printf("Processed %v\n", file)

	gz.Close()
	f.Close()

	wg.Done()
}

func processOutput(rows <-chan []string, wg *sync.WaitGroup) {
	i := 0

	for row := range rows {
		write(row[0], row[1:])
		i++

		if i%500000 == 0 {
			log.Printf("Processed %v\n", i)
		}
	}

	wg.Done()
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

	bufSize := 10000
	rowsGroup := []chan []string{
		make(chan []string, bufSize),
		make(chan []string, bufSize),
		make(chan []string, bufSize),
		make(chan []string, bufSize),
		make(chan []string, bufSize),
		make(chan []string, bufSize),
		make(chan []string, bufSize),
		make(chan []string, bufSize),
		make(chan []string, bufSize),
		make(chan []string, bufSize),
	}

	var wgR sync.WaitGroup

	for _, file := range files {
		wgR.Add(1)

		go processFile(*destDir, file, &wgR, rowsGroup)
	}

	var wgW sync.WaitGroup

	for _, rows := range rowsGroup {
		wgW.Add(1)

		go processOutput(rows, &wgW)
	}

	wgR.Wait()
	log.Printf("Finished reading\n")

	for _, rows := range rowsGroup {
		close(rows)
	}

	wgW.Wait()
	log.Printf("Finished writing\n")
}
