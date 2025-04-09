# NATS JetStream Cluster Testing

## Cluster management
* Start: `make start`
* Stop: `make stop`
* Logs: `make logs`
* List containers: `make ps`
* Kill random NATS node: `make kill`
* Explicitly kill JetStream leader node: `make kill_leader`
* Stop and fully cleanup JetStream data volumes: `make destroy`

## Get info
`make js_status`

`make account_info`

## Run nats test clients
`make test.py`

`make test.go`

## Example
```shell
# Start
$ make up
...

# Run test
$ make test.py
Connected to ParseResult(scheme='nats', netloc='a:a@localhost:4222', path='', params='', query='', fragment='')
Connect urls: ['172.22.0.2:4222', '172.22.0.3:4222', '172.22.0.4:4222']
StreamState(messages=0, bytes=0, first_seq=0, last_seq=0, consumer_count=0, deleted=None, num_deleted=None, lost=None, subjects=None)
Received a message on foo: Message 1
Received a message on foo: Message 2
Received a message on foo: Message 3

# Run again, notice number of messages in stream
$ make test.py
Connected to ParseResult(scheme='nats', netloc='a:a@localhost:4222', path='', params='', query='', fragment='')
Connect urls: ['172.22.0.3:4222', '172.22.0.2:4222', '172.22.0.4:4222']
StreamState(messages=3, bytes=126, first_seq=1, last_seq=3, consumer_count=1, deleted=None, num_deleted=None, lost=None, subjects=None)
Received a message on foo: Message 1
Received a message on foo: Message 2
Received a message on foo: Message 3

# Get status, notice persisted stream files on each node
$ make js_status
...
-rw-------    1 root     root           252 Apr  9 16:34 1.blk
-rw-------    1 root     root           252 Apr  9 16:34 2.blk
-rw-------    1 root     root           252 Apr  9 16:34 2.blk

# Stop a node
$ make kill
nats-nats2-1

# Test again, note we lost a node reporting back in the server's INFO message
$ make test.py
Connected to ParseResult(scheme='nats', netloc='a:a@localhost:4222', path='', params='', query='', fragment='')
Connect urls: ['172.22.0.3:4222', '172.22.0.2:4222']
StreamState(messages=6, bytes=252, first_seq=1, last_seq=6, consumer_count=1, deleted=None, num_deleted=None, lost=None, subjects=None)
Received a message on foo: Message 1
Received a message on foo: Message 2
Received a message on foo: Message 3

# Get status, node down
$ make js_status
...
-rw-------    1 root     root           378 Apr  9 16:37 1.blk
-rw-------    1 root     root            64 Apr  9 16:36 index.db
Error response from daemon: container a6214e486e163fbc58a48cfa5a784a3bbfe86587edfc6eae7a2776b10a3dd662 is not running
-rw-------    1 root     root           378 Apr  9 16:37 2.blk
-rw-------    1 root     root            64 Apr  9 16:36 index.db

# Kill another node
$ make kill
nats-nats3-1

# Can no longer use JetStream
$ make test.py
Connected to ParseResult(scheme='nats', netloc='a:a@localhost:4222', path='', params='', query='', fragment='')
Connect urls: ['172.22.0.2:4222']
Traceback (most recent call last):
...
nats.js.errors.ServiceUnavailableError: nats: ServiceUnavailableError: code=503 err_code=10008 description='JetStream system temporarily unavailable'

# Restart
$ make start
...

# Test again; messages retained
$ make test.py
Connected to ParseResult(scheme='nats', netloc='a:a@localhost:4222', path='', params='', query='', fragment='')
Connect urls: ['172.22.0.4:4222', '172.22.0.2:4222', '172.22.0.5:4222']
StreamState(messages=12, bytes=504, first_seq=1, last_seq=12, consumer_count=1, deleted=None, num_deleted=None, lost=None, subjects=None)
Received a message on foo: Message 1
Received a message on foo: Message 2
Received a message on foo: Message 3

# Specifically kill JetStream cluster leader
$ make kill_leader 
nats-nats2-1

$ make test.py
Connected to ParseResult(scheme='nats', netloc='a:a@localhost:4222', path='', params='', query='', fragment='')
Connect urls: ['172.22.0.5:4222', '172.22.0.3:4222']
StreamState(messages=24, bytes=1008, first_seq=1, last_seq=24, consumer_count=1, deleted=None, num_deleted=None, lost=None, subjects=None)
Received a message on foo: Message 1
Received a message on foo: Message 2
Received a message on foo: Message 3
$ make test.py
Connected to ParseResult(scheme='nats', netloc='a:a@localhost:4222', path='', params='', query='', fragment='')
Connect urls: ['172.22.0.3:4222', '172.22.0.5:4222']
StreamState(messages=27, bytes=1134, first_seq=1, last_seq=27, consumer_count=1, deleted=None, num_deleted=None, lost=None, subjects=None)
Received a message on foo: Message 1
Received a message on foo: Message 2
Received a message on foo: Message 3
```