# DistrKV

DistrKV is a distributed Key-Value Store

### Key Features

- GET/SET requests to store and retrieve key-value pairs
- Multiple shards for storing key-value pairs
- Keys are hashed and then shard index is determined based upon the hash and shard count using static sharding.

### Interacting with the Store

- Start the store using the ```start.sh``` script
- Make curl requests of the form ```curl localhost:<port>/get?key=<key>``` to make get requests
- Make curl requests of the form ```curl localhost:<port>/set?key=<key>&value=<value>``` to make set requests
