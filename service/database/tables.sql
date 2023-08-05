CREATE TABLE transactions (
    Id INT AUTO_INCREMENT PRIMARY KEY,
    Reference INT UNIQUE,
    SourceAccount VARCHAR(100),
    DestinationAccount VARCHAR(100),
    Amount INT,
    Date TIMESTAMP,
    SourceBank VARCHAR(100),
    DestinationBank VARCHAR(100),
    Type VARCHAR(50) CHECK IN ("atm", "bank", "app"),
);