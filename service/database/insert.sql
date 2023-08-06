INSERT INTO transactions (Reference, SourceAccount, DestinationAccount, Amount, Date, SourceBank, DestinationBank, Type) 
    VALUES (%s, %s, %s, %s, %s, %s, %s, %s);