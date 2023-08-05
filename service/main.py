import json
import mysql.connector



# open and read config file
f = open('config.json')
config = json.load(f)

# closing file
f.close()

# display configs
print(json.dumps(config, indent=4))

# open database connection
connection = mysql.connector.connect(
  host=config['DB_HOST'],
  port=config['DB_PORT'],
  user=config['DB_USER'],
  password=config['DB_PASS'],
  database=config['DB_NAME']
)

# open a cursor
cursor = connection.cursor()

# check tables
cursor.execute("SHOW TABLES")

for x in cursor:
  print(x)
