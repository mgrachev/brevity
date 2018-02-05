# Brevity
Another useless service for shortening links :trollface:

## Install

Clone the repository:

```bash
$ git clone https://github.com/mgrachev/brevity
$ cd brevity
```

Install dependencies:

```bash
$ go get github.com/rnubel/pgmgr
$ go get -u github.com/kardianos/govendor
$ govendor sync
```

Create a database and run migrations:

```bash
$ pgmgr db create
$ pgmgr db migrate
```

Run the HTTP-server:

```bash
$ PG_CONNECTION_URL=postgres://... go run cmd/brevity-http-server/brevity-http-server.go
```
