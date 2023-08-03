# Fraud Detection

Financial transactions fraud detection using Machine Learning. Implementing
a model for detecting the frauds in financial transactions in a banking system.
In this project we are going to create a dataset from a bank system database. After that,
we are going to make ```SQL``` queries in order to get transactions from bank database.

Furthermore, we are going to build a machine learning model to detect the fraud transactions and
give us a report on them.

## Input data

Each transaction comes with a data as bellow:

```json
{
  "id": 10,
  "reference" : 10399301901,
  "source account": "1910391i03",
  "dest account": "193849018029840",
  "amount" : 100,
  "date": "Jan 12902",
  "source bank": "Meli",
  "dest bank": "Some place",
  "type": "ATM/Bank/App"
}
```

## Output

By getting a list of transactions we need to return the ones that are possibly a fraud.

```json
[
  {
    "transactions id": 10,
    "fruad type": "robbery/fake/etc" 
  }
]
```
