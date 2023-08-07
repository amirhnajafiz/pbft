from model import Transatcion



class Provider(object):
    def __init__(self, batch):
        self.batch = batch
        
        # get insert query
        with open('./database/insert.sql', 'r') as file:
            self.query = file.read()
    
    def generate(self, db):
        # create a batch list
        list = []
        for _ in range(0, self.batch):
            list.append(Transatcion().list())
            
        cursor = db.cursor()
        
        # insert many
        cursor.executemany(self.query, list)
        
        print(cursor.rowcount, "record inserted.")
        
        db.commit()
        
        cursor.close()
        db.close()
