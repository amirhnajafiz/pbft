class Provider(object):
    def __init__(self, batch):
        self.batch = batch
        # get insert query
        with open('./database/insert.sql', 'r') as file:
            self.query = file.read()
    
    def generate(self, db):
        pass