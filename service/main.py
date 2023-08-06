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

# make migration
if config['MIGRATE']:
  # open a cursor
  cursor = connection.cursor()
  with open('database/tables.sql', 'r') as file:
      query = file.read()
      cursor.execute(query)

  print("Migrated!")
  cursor.close()


# create provider
from provider import Provider

p = Provider(config['BATCH'])
p.generate(connection)
