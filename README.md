# go-redis-from-scratch

A Redis-compatible in-memory server built from scratch in Go — no external dependencies.  
It implements the RESP (Redis Serialization Protocol) wire format, meaning any standard Redis client (`redis-cli`, `ioredis`, `jedis`, etc.) can connect and talk to it out of the box.

> Built as a learning project to understand how Redis works under the hood — TCP servers, RESP parsing, TTL expiry strategies, pub/sub messaging, and concurrency in Go.

---

## Prerequisites

- [Go](https://golang.org/dl/) 1.18 or higher

---

## Installation

```bash
git clone https://github.com/your-username/go-redis-from-scratch.git
cd go-redis-from-scratch
```

No external dependencies — the standard library is all that's used.

---

## Running

**Development (run directly):**
```bash
go run main.go
```

**Production (build first):**
```bash
go build -o go-redis .
./go-redis
```

Server listens on `127.0.0.1:6379` by default. Override the bind address with the `REDIS_BIND` environment variable:

```bash
REDIS_BIND=0.0.0.0:6380 go run main.go
```

---

## Connecting

Use any Redis client or `redis-cli`:

```bash
redis-cli -p 6379
```

Example session:

```
127.0.0.1:6379> PING
PONG
127.0.0.1:6379> SET name "alice"
OK
127.0.0.1:6379> GET name
"alice"
127.0.0.1:6379> SETEX session 60 "token123"
OK
127.0.0.1:6379> TTL session
(integer) 59
127.0.0.1:6379> INCR counter
(integer) 1
127.0.0.1:6379> LPUSH queue "job1" "job2"
(integer) 2
127.0.0.1:6379> LRANGE queue 0 -1
1) "job2"
2) "job1"
```

---

## Supported Commands

| Category  | Commands                                    |
|-----------|---------------------------------------------|
| Core      | `SET`, `GET`, `DEL`, `EXISTS`, `PING`       |
| Expiry    | `EXPIRE`, `TTL`, `SETEX`                    |
| Counters  | `INCR`, `DECR`, `INCRBY`                    |
| Lists     | `LPUSH`, `RPUSH`, `LPOP`, `RPOP`, `LRANGE` |
| Pub/Sub   | `SUBSCRIBE`, `PUBLISH`                      |

---

## How It Works

- **Lazy expiry** — expired keys are detected and deleted on `GET`
- **Active expiry** — a background goroutine sweeps and cleans up expired keys every 10 seconds
- **Concurrency** — `sync.RWMutex` protects all reads and writes to the store and pub/sub maps
- **Pub/Sub** — subscribers hold open TCP connections; `PUBLISH` writes the message directly to each live connection
- **RESP protocol** — commands are parsed from the raw TCP stream following the Redis wire format (`*`, `$` prefixes)

---

## Project Structure

```
go-redis-from-scratch/
├── main.go              # Entry point — boots store, pubsub, and TCP server
├── go.mod
├── server/
│   └── tcp_server.go    # Accepts connections, spawns a goroutine per client
├── protocol/
│   └── resp.go          # RESP parser and command dispatcher
└── storage/
    ├── cache.go         # In-memory key/value store with TTL and list support
    └── pubsub.go        # Pub/Sub engine backed by raw TCP connections
```

---

## License

MIT
