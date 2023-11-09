# tickets-api
Tickets API Simulation

## Seed the Database

* Install dependencies
  * `pip install -r requirements.txt`
* Update the `config/config.json` with the DB information
* Execute the Python script
  * Secure CRDB: `python db-seed/seed.py --secure --config config/config.json --num_users 10000 --num_purchases 10000 --num_events 10000 --num_cancellations 10000 --num_payments 10000`
  * Insecure CRDB: `python db-seed/seed.py --config config/config.json --num_users 10000 --num_purchases 10000 --num_events 10000 --num_cancellations 10000 --num_payments 10000`

## Tickets API Endpoints

* GET: `/user/:uuid/purchases` - get user purchases
* GET: `/user/:uuid/purchases/cancellations` - get user cancellations
* GET: `/search/user/:uuid` - search users by ID `http://localhost:3001/search/user/27912e48-88f5-4dca-b549-88da7bdda1f4` for go-pg
* GET: `/search/user/:uuid` - search users by ID `http://localhost:3002/search/user/27912e48-88f5-4dca-b549-88da7bdda1f4` for pgx

## Start GO-PG Tickets API Endpoint

```shell
cd gopg-api

go mod init gopg-api
go mod tidy
go run *.go
```

### Collect Test Data Points for the k6 Test

```sql
SELECT purchase.user_id
  FROM purchases AS purchase LEFT
    JOIN users AS users
     ON users.id = purchase.user_id LEFT
    JOIN events AS events
     ON events.id = purchase.event_id order by random() limit 100 ;
```

The results will need to be added to the k6 `const names = []`

### Run K6 Stress Test

```
k6 run k6-gopg.js
```

## Start PGX Tickets API Endpoint

```shell
cd pgx-api

go mod init pgx-api
go mod tidy
go run *.go
```
### Run K6 Stress Test

```shell
k6 run k6-pgx.js
```

## Implicit/Explicit Transaction Example

* This example focuses on querying the `users` table
* This example requires Jaeger (can be deployed in a docker container) to trace the API execution and it's SQL statement basic metrics
* This example is built to test Read Committed transactions in CockroachDB 23.2.x (beta at the moment of this writing)

```shell
cd read-commit

go mod init implicit-explicit
go mod tidy
go run *.go
```

**NOTE:** for K6 Stress test, the ramp up in this examples have been set to 1m.

### Test Implicit Transactions

* Slightly modify your K6 script to hard code some user UUIDs in an array
  * `const uuids = [...uuids...]`
* API endpoint: http://localhost:8080/implicit/users/$uuid
* Run the stress test
  * `k6 run --vus 500 --duration 5m k6-implicit.js`
  
### Test Explicit Transactions (v23.2.x+)

* Slightly modify your K6 script to hard code some user UUIDs in an array
  * `const uuids = [...uuids...]`
* API endpoint: http://localhost:8080/explicit/users/$uuid
* Enable Read Commit in CockroachDB
  * `root@192.168.86.74:26257/defaultdb ?> SET CLUSTER SETTING sql.txn.read_committed_syntax.enabled = 'true'; `
  * At the session level set: `set default_transaction_isolation = 'read committed';`
* Run the stress test
  * `k6 run --vus 500 --duration 5m k6-explicit.js`

