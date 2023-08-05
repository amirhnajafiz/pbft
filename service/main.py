import json



# open and read config file
f = open('config.json')
config = json.load(f)

# closing file
f.close()

# display configs
print(json.dumps(config, indent=4))
