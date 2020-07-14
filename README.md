###Bid placement app

####Setup with docker-compose
* Install docker and docker-compose
* Run `docker-compose up`

####Setup without docker
* Install golang and postgres 
* Connect to postgres and run `create database bids;`
* Execute scripts `./migrations/*.up.sql` in postgres
* Run `go run ./main.go` or `make run` to start application

By default app is deployed on `localhost:8080` and will connect to postgres with host `postgres` and default credentials.

####Testing
* To test the app execute http post request with body from the task:
`curl -X POST -H "Source-Type: game" -d '{"state": "win", "amount": "11.11", "bidId": "12345"}' http://localhost:8080/bids`
* Or run integration tests with `go test -v -p 1 ./...` or `make test`
* There are only integration tests so the database has to be started.
