# Bitcoin ancestries

This codebase will produce some stats on the ancestry of each transaction of the Bitcoin network.

In the Bitcoin network each transaction can have multiple inputs. Each input is another
transaction. This way each coin you spent is traceable all the way down to when
it was minted.

Coinbase transactions, i.e. transactions that create bitcoin don't have any
inputs. This means that each transaction has an ancestral tree where leaf nodes
are coinbase transactions.

For example, we can visualise one such heritage tree as follows

```.
transaction
├── transaction
│   ├── coinbase
│   ├── coinbase
├── trancsaction
│   ├── coinbase
│   └── coinbase
└── transaction
    ├── transaction
        ├── trasaction
            ├── transaction
                └──coinbase
```


This yields the following interesting question:

* What has been the longest chain ever observed in bitcoin (i.e. find a transaction that forms a tree of a maximum height)

 # Answer

Up to block 607393 (minted 2019-12-09 17:08:09) the longest chain is of length 2886138.

The head of the tree is the transaction [8533dc23fc315f91b5fcd7dc931a3fabb7ac999cca108bd79a7a081a1bd95ae9][blockchain], included in block 607392.

It is 6 hops away from the closest coinbase (from block 607275) transaction and 2886138 hops away from the most distant one (from block 9).

# Another interesting result

Transaction [8d9e9fe8fd08af344763275521e2010bb7cf9eb042d90157011178e04bd3067f][blockchain2] is
the head of the tree with the longest shortest path to a coinbase transaction.
It is 25348 hops away from a coinbase.

# Methodology

Following Golang code was written to process the entire blockchain and to find the height of each tree and more.

* [cmd/blocksplitter][splitter] groups raw BigQuery dump of all transactions by block
* [cmd/blocksorter][sorter] applies [topological sorting](https://en.wikipedia.org/wiki/Topological_sorting) to each block
  (now in its own file) making sure that transactions in each file appears in
  a chronoligical order
* [cmd/blockexplorer][explorer] computes ancestry characteristics (hops_min,
  hops_max, mined_block_min, mined_block_max) for each transaction.

# Notes

This is research code showing how to process large datasets relatively quickly.
My laptop (Intel(R) Core(TM) i7-8565U CPU @ 1.80GHz) processes >100 gb of data
within a few hours (>1.1 billion rows).

Due to a bug in Bitcoin some transaction hashes were repeated. However those
hashes don't break the code since they are both coinbase. This bug was [fixed][bip]
and duplicate hashes are no longer observed.

# Source of data

BigQuery's public dataset:


```.sql
SELECT t.block_number, t.hash, i.value, i.spent_transaction_hash FROM `bigquery-public-data.crypto_bitcoin.transactions` t LEFT JOIN UNNEST(t.inputs) i
```

[blockchain]: https://www.blockchain.com/es/btc/tx/8533dc23fc315f91b5fcd7dc931a3fabb7ac999cca108bd79a7a081a1bd95ae9
[blockchain2]: https://www.blockchain.com/es/btc/tx/8d9e9fe8fd08af344763275521e2010bb7cf9eb042d90157011178e04bd3067f
[bip]: https://github.com/bitcoin/bips/blob/master/bip-0034.mediawiki
[splitter]: https://github.com/afiodorov/bitcoin_ancestries/tree/master/cmd/blocksplitter
[sorter]: https://github.com/afiodorov/bitcoin_ancestries/tree/master/cmd/blocksorter
[explorer]: https://github.com/afiodorov/bitcoin_ancestries/tree/master/cmd/blockexplorer
