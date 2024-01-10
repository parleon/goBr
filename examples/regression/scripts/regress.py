import pandas as pd
import numpy as np
import statsmodels.api as sm
import sys


def parse_path(path):
    return path.split('.')

def strip_columns(path):
    parsed = parse_path(path)
    y_df = pd.read_csv("./scripts/datasets/"+parsed[0]+".csv")
    y_df = y_df.loc[:, y_df.columns.intersection(['Date', parsed[1]])]
    y_df.set_index("Date")
    y_df = y_df.rename(columns={c: parsed[0]+"."+c for c in y_df.columns if c != 'Date'})
    y_df = y_df.dropna(axis=0)
    y_df = y_df.dropna(axis=1)
    y_df = y_df.copy()
    return y_df

def create_regression_table(y, x):
    y = strip_columns(y)

    for i in range(len(x)):
        x[i] = strip_columns(x[i])

    for i in range(len(x)):
        y = y.merge(x[i], on='Date', how='inner')

    y = y.drop('Date', axis=1)
    y = y.dropna(axis=1)
    y = y.copy()

    y['Period'] = range(1,len(y)+1)
    y['Period_Squared'] = np.square(y['Period'])
    y['Period_log'] = np.log(y['Period'])

    return y

def perform_regression(regression_table):


    y = regression_table[regression_table.keys()[0]]
    X = regression_table[regression_table.keys()[1:]]

    X = sm.add_constant(X)

    model = sm.OLS(y, X).fit()

    return model


def main():
    y = sys.argv[1]
    x = sys.argv[2:]

    regression_table = create_regression_table(y, x)

    print(perform_regression(regression_table).summary())


if __name__ == '__main__':
    main()