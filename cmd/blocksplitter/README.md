# What is it?

This is a way to split all transactions downloaded from BigQuery in the bitcoin
chain into separate blocks

# What is the output like?

```.bash
> head -n5 531351.csv

hash,value,spent_transaction_hash
6d808e3bca5bc70e64770a59fbe73651c7d13389b7c60344672d271f8e108a9b,1182159,54c7c19d72b1c00824c2a2f4835b2c90286eb7d3f56eb26b6935eb4fd8690f89
b398eb16418462726bec90192688ba6986a9c3a5ce2fd146b59589f08bcd3b2c,60959032,7876b5d19f6f6a89ac5a216bcd4ff57c3aa45f1df8f59e3d2ab07463b7b53a81
303c7d04b03637045cb72e38da9dd2b6abb272a061b0c9a47d14de79304e086b,559231,c073f32fab1126798b34eaa48472d27f8b7c7bfdf54e6737b73eeaed350d6305
ced69e1fe5708e0a022b61357d482bc766928a640e16771b1ceaf88542bb2db0,679604776,ea4ff9ec5b777c7559313a26ff0ea6be11855fbad61eeba49657f157c83272ea
```

# How to use it?

Run

```.sql
SELECT t.block_number, t.hash, i.value, i.spent_transaction_hash FROM `bigquery-public-data.crypto_bitcoin.transactions` t LEFT JOIN UNNEST(t.inputs) i
```

& save the output to a table & export that table to a bucket (gzipped for a speedy download)

Then

```.bash
gsutil -m cp -n gs://tmp/bitcoin* ~/data/bitcoin
```

We now have data of the form

```.bash
block_number,hash,value,spent_transaction_hash
531351,6d808e3bca5bc70e64770a59fbe73651c7d13389b7c60344672d271f8e108a9b,1182159,54c7c19d72b1c00824c2a2f4835b2c90286eb7d3f56eb26b6935eb4fd8690f89
531351,b398eb16418462726bec90192688ba6986a9c3a5ce2fd146b59589f08bcd3b2c,60959032,7876b5d19f6f6a89ac5a216bcd4ff57c3aa45f1df8f59e3d2ab07463b7b53a81
```

Finally run

```.bash
go run main.go -s ~/data/bitcoin -d ~/data/blocks
```

and it will create csv files for each block in the destination folder containing data.

# QA

Check that we have correct numbers against each block

```.bash
find ~/data/blocks -name "*.csv" -print0 | xargs -0 -i sh -c 'v="{}" && echo -n "$v " && tail -n+2 $v | wc -l' > /tmp/actual_sizes
```

BigQuery:

```.sql
SELECT COUNT(*) counts, block_number FROM `bigquery-public-data.crypto_bitcoin.transactions` t LEFT JOIN UNNEST(t.inputs) i GROUP BY block_number ORDER BY block_number
```

QA:

```.python
import pandas as pd

actual = pd.read_csv('/tmp/actual_sizes', header=None, names=['block_number', 'counts'], sep=" ")
actual['block_number'] = actual['block_number'].apply(lambda x: x.split('/')[-1].split('.')[0]).astype(int)
actual.sort_values('block_number', inplace=True)
actual = actual[['counts', 'block_number']].reset_index(drop=True)

expected = pd.read_csv('~/data/qa.csv')
expected = expected[expected['block_number'] <= actual['block_number'].max()]

pd.testing.assert_frame_equal(actual, expected)
```

Curiously I got 1 missing line in one file. After manually extracting that block counts match up perfectly now.
