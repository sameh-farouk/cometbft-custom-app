# CometBFT Custom App with Multi-Database Support

This is a custom blockchain application built on top of CometBFT with support for multiple database backends.

## Supported Databases

- **Badger DB** (default): A fast key-value database
- **Pebble DB**: A performance-focused key-value store used in CockroachDB
- **TigerBeetle DB**: A specialized financial database designed for financial transactions

## Quick Start

```bash
git clone https://github.com/sameh-farouk/cometbft-custom-app.git
cd cometbft-custom-app

# Build with default database (Badger, pebble)
make build-linux
DB_TYPE=badger make start-localnet
# or
DB_TYPE=pebble make start-localnet
make apply-latency
```

## Using Different Database Backends

You can specify which database to use with the `DB_TYPE` environment variable:

```bash
# Build with TigerBeetle DB support
make build-tigerbeetle
DB_TYPE=tigerbeetle make start-localnet
```

## Example Usage

```bash
# Send a transaction
curl -s 'localhost:26657/broadcast_tx_commit?tx="1=2=50=93038a25e2806ea16771281f734b8dc28600082ef05f59e39a6765c5f73dd115991dd91e5c4af11e7d20b5170519eee281d873dcb625730e95d8e0fb3ff23902"'

# Query account balances
curl -s 'localhost:26657/abci_query?data="1"'
echo OTk5OTk5OTUw | base64 --decode  # 999999950

curl -s 'localhost:26657/abci_query?data="2"'
echo MTAwMDAwMDA1MA== | base64 --decode  # 1000000050

# send batched transactions
curl -s 'localhost:26657/broadcast_tx_commit?tx="1=2=50=93038a25e2806ea16771281f734b8dc28600082ef05f59e39a6765
c5f73dd115991dd91e5c4af11e7d20b5170519eee281d873dcb625730e95d8e0fb3ff23902:1=2=50=93038a25e2806ea16771281f734b8dc28600082ef05f59e39a6765c5f73dd115991dd91e5c4af11e7d20b5170519eee281d873dcb625730e95d8e0fb3ff23902:1=2=50=93038a25e2806ea16771281f734b8dc28600082ef05f59e39a6765c5f73dd115991dd91e5c4af11e7d20b5170519eee281d873dcb625730e95d8e0fb3ff23902"'
```

## Database Configuration

Each database can be configured with additional options:

- **Badger DB**: Use `-db-path` to specify the database directory
- **Pebble DB**: Use `-db-path` to specify the database directory
- **TigerBeetle DB**: 
  - Use `-tb-addresses` to specify TigerBeetle server addresses (comma-separated)
  - Use `-tb-cluster-id` to specify the TigerBeetle cluster ID
