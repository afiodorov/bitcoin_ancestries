package main

import (
	"log"
	"path"
	"path/filepath"
	"strconv"
)

func withoutExt(filename string) string {
	filename = path.Base(filename)
	extension := filepath.Ext(filename)
	name := filename[0 : len(filename)-len(extension)]

	return name
}

func getBlockNumber(file string) int {
	bn, err := strconv.Atoi(withoutExt(file))
	if err != nil {
		log.Fatal(err)
	}

	return bn
}

func uniqueOrdered(list []string) []string {
	seen := make(map[string]struct{})
	res := make([]string, 0)

	for _, e := range list {
		if _, ok := seen[e]; ok {
			continue
		}

		seen[e] = struct{}{}

		res = append(res, e)
	}

	return res
}

func unique(list []string) []string {
	seen := make(map[string]struct{})

	for _, e := range list {
		seen[e] = struct{}{}
	}

	res := make([]string, 0, len(seen))
	for k := range seen {
		res = append(res, k)
	}
	return res
}
