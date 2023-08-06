from model import Transatcion



class Provider(object):
    def __init__(self, batch):
        self.batch = batch
        
        # get insert query
        with open('./database/insert.sql', 'r') as file:
            self.query = file.read()
    
    def generate(self, db):
        for _ in range(0, self.batch):
            cursor = db.cursor()

            cursor.execute(self.query, Transatcion().list())
            db.commit()
            
            print(cursor.rowcount, "record inserted.")
            
            cursor.close()
