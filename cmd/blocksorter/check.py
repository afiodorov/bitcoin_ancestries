#!/usr/bin/env python
import argparse
import sys
from pathlib import Path

import numpy as np
import pandas as pd


def is_correct(df):
    for i in range(len(df)):
        intersect = np.intersect1d(df.iloc[i, -1], df.iloc[i+1:, 0])
        if len(intersect) > 0:
            return False, intersect[0]
    return True, ""


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='QA sorter')
    parser.add_argument('--source', type=str, help='Source directory')
    args = parser.parse_args()

    for f in Path(args.source).glob("*.csv"):
        df = pd.read_csv(f).astype(str)
        correct, failed = is_correct(df)
        if not correct:
            print(f"File {f} FAILED. Problem: {failed}")
            sys.exit(1)
        print(f"File {f} PASSED")
