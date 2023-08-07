INSERT INTO transactions (Reference, SourceAccount, DestinationAccount, Amount, CreatedAt, SourceBank, DestinationBank, Type) 
    VALUES (%s, %s, %s, %s, %s, %s, %s, %s);