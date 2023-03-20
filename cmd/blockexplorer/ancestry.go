package main

import "log"

type Ancestry struct {
	Hash          string `pg:",notnull"`
	BlockNumber   int
	HopsMin       int
	HopsMax       int
	MinedBlockMin int
	MinedBlockMax int
}

func makeAncestry(transaction string, children []string, blockNumber int, resHolder ResHolder) Ancestry {
	out := Ancestry{
		Hash:        transaction,
		BlockNumber: blockNumber,
	}

	isCoinbase := len(children) == 1 && children[0] == ""
	if isCoinbase {
		out.MinedBlockMin = blockNumber
		out.MinedBlockMax = blockNumber
		return out
	}

	if len(children) == 0 {
		log.Fatalf("No children for transaction %v\n", transaction)
	}

	firstChildAncestry, ok := resHolder.LookUp(children[0])
	if !ok {
		log.Fatalf("Input %v of transaction %v not found\n", children[0], transaction)
	}

	hopsMin := firstChildAncestry.HopsMin
	hopsMax := firstChildAncestry.HopsMax
	blockMin := firstChildAncestry.MinedBlockMin
	blockMax := firstChildAncestry.MinedBlockMax

	for _, child := range children[1:] {
		childAncestry, ok := resHolder.LookUp(child)
		if !ok {
			log.Fatalf("Input %v of transaction %v not found\n", child, transaction)
		}

		if childAncestry.HopsMin < hopsMin {
			hopsMin = childAncestry.HopsMin
		}
		if childAncestry.HopsMax > hopsMax {
			hopsMax = childAncestry.HopsMax
		}
		if childAncestry.MinedBlockMin < blockMin {
			blockMin = childAncestry.MinedBlockMin
		}
		if childAncestry.MinedBlockMax > blockMax {
			blockMax = childAncestry.MinedBlockMax
		}
	}

	out.HopsMin = hopsMin + 1
	out.HopsMax = hopsMax + 1
	out.MinedBlockMin = blockMin
	out.MinedBlockMax = blockMax
	return out
}
