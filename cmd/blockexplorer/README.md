# What is it?

This actually computes stats on the ancestry of each transaction of the Bitcoin network.

The following stats stats are computed:

* hops_min - minimum number of hops required to reach a coinbase transaction
* mined_block_min - the block number where the above coinbase was included
* hops_max - height of the ancestry tree
* mined_block_min - the block number where the earliest coinbase was included


# How to run?

Start a postgres instance:

```.bash
docker run -d -p 5433:5432 --name blockexplorer -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=writer -e POSTGRES_DB=pgdatabase postgres:11.7
```

No need to create a schema - the binary will create it. Make sure you have run
`blocksplitter` and `blocksorter` first.

Then do

```.bash
go build
./blockexplorer -s sorted_blocks_folder
```

The computation is resumable via a `-f block_to_resume_from` flag.

# Example logs

```.bash
2019/12/16 22:35:38 Block 437424. Inserting 501549 transactions...
2019/12/16 22:37:54 Block 437740. Inserting 501915 transactions...
2019/12/16 22:39:11 Block 438007. Inserting 500212 transactions...
2019/12/16 22:40:25 Block 438298. Inserting 500709 transactions...
2019/12/16 22:41:38 Block 438639. Inserting 501899 transactions...
2019/12/16 22:42:50 Block 438960. Inserting 502480 transactions...
2019/12/16 22:44:02 Block 439251. Inserting 500343 transactions...
2019/12/16 22:45:16 Block 439535. Inserting 502181 transactions...
2019/12/16 22:46:31 Block 439809. Inserting 500853 transactions...
2019/12/16 22:47:39 Block 440051. Inserting 500726 transactions...
```

takes 0.274 seconds per block ± 0.06 for later blocks.

# Methodology

The computation uses postgres to store ancestries of each transaction. When we
load a block we can instantly see which previous transactions are needed to
process such block. We then load previous transactions from the db (there is
also a LRU cache to hit the database less frequently) then write to the db
every time we process `syncStep = 100000` transactions.

During the experimentations I noticed that using a hash index on the
transaction id speeds up the database interactions by quite a lot.
