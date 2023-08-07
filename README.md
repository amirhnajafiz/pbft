# Fraud Detection

Financial transactions fraud detection using Machine Learning. Implementing
a model for detecting the frauds in financial transactions in a banking system.
In this project we are going to create a dataset from a bank system database. After that,
we are going to make ```SQL``` queries in order to get transactions from bank database.

Furthermore, we are going to build a machine learning model to detect the fraud transactions and
give us a report on them.

## Input data

Each transaction comes with a data structure in our system database as bellow:

```sql
CREATE TABLE transactions (
    Id INT AUTO_INCREMENT PRIMARY KEY,
    Reference INT UNIQUE,
    SourceAccount VARCHAR(100),
    DestinationAccount VARCHAR(100),
    Amount INT,
    CreatedAt VARCHAR(100),
    SourceBank VARCHAR(100),
    DestinationBank VARCHAR(100),
    Type VARCHAR(50),
    CHECK (Type IN ("atm", "bank", "app"))
);
```

Transactions have a unique reference, source and dest account and banks, the amount of money for that transaction,
the time of transaction and finally transaction type which is from ```ATM```, ```Bank```, or ```App```.

## Output

By getting a list of transactions we need to return the ones that are possibly fraud ones. We don't
need to specify the reason, we just need to find the fraud ones in order to get checked by the system
administrators.

```json
[
  {
    "transactions id": 10,
    "fruad type": true
  }
]
```
