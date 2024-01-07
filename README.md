# E-Commerce Example

## Description

This is a practice project using go, nats, eventstoreDb.

## Setup

## Run

```
# run the nats & esdb services
docker-compose up -d

# run the gateway
go run cmd/gateway/main.go

# run the command subscribers
go run cmd/shopping-cart-command-handler/main.go

# example use case: open a shopping cart
curl -XPOST localhost:3000/shopping-cart
```

## Todos

- Add missing features
- Clean inconsistencies in naming and such
- Look into nats jetstream as replacement for esdb
- Testing
- Documentation
- Experiment with projections and snapshots
