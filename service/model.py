from data import Data


class Transatcion(object):
    def __init__(self):
        d = Data()
        
        self.reference = d.reference()
        self.source = d.bank_account()
        self.dest = d.bank_account()
        self.amount = d.amount()
        self.date = d.date()
        self.sbank = d.bank()
        self.dbank = d.bank()
        self.rtype = d.rtype()
    
    def list(self):
        return [
            self.reference,
            self.source,
            self.dest,
            self.amount,
            self.date,
            self.sbank,
            self.dbank,
            self.rtype
        ]
