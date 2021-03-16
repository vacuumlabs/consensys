
# VAX

This is a monorepo that contains all microservies for the VAX distributed application written in the Go language.

Use Docker Compose to start all the microservicse `docker-compose up`, or `docker-compose up --build` to force re-build the Docker images.

Services are built via a multi-stage Dockerfile, so first build is slow because it downloads the dependencies. Subsequent builds are faster thanks to Docker layer caching.  

Communication between services is done via Redis Pub/Sub messaging, instead of Kafka.

We have one manufacturer "manu_01", one authority "auth_01", one customer "cust_01".  These are defined in the `docker-compose.yml` file via environment variables.

All data is stored in-memory only. 

## Postman

Use the `VAX.postman_collection.json` file.

## cURL examples
Or use cURL for testing.

Create a new shipment, that will be sent to an authority and to a customer:
```
curl -v -X POST -H "Content-Type: application/json" \
     -d '{"vaccineName": "vacc 01", "quantity": 1000, "expirationDays": 90, "authorityId": "auth_01", "customerId": "cust_01"}' \
     http://localhost:8082/api/v1/shipment
```

Get information about an event from the notary service via its id:
```
curl -v http://localhost:8080/api/v1/event/:event_id
```

Create a new event for a hash:
```
curl -v -X POST -H "Content-Type: application/json" \
     -d '{"hash": "...HASH..."}' \
     http://localhost:8080/api/v1/event
```

## Unit tests

Only one unit test currently, for testing thread safety of the nonce provider. Run it via `go test` in the `./pkg/memory` folder. GO SDK must be installed for this step. The test can be also run with Go's race detector to verify that there are no race conditions `go test -race`.
