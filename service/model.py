class Transatcion(object):
    def __init__(self, reference, source, dest, amount, date, sbank, dbank, type):
        self.reference = reference
        self.source = source
        self.dest = dest
        self.amount = amount
        self.date = date
        self.sbank = sbank
        self.dbank = dbank
        self.type = type
    
    def list(self):
        return [
            self.reference,
            self.source,
            self.dest,
            self.amount,
            self.date,
            self.sbank,
            self.dbank,
            self.type
        ]
