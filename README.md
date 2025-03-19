# CometBFT Custom App with Multi-Database Support

This is a custom blockchain application built on top of CometBFT with support for multiple database backends.

## Supported Databases

- **Badger DB** (default): A fast key-value database
- **Pebble DB**: A performance-focused key-value store used in CockroachDB
- **TigerBeetle DB**: A specialized financial database designed for financial transactions (WIP)

## Quick Start

```bash
git clone https://github.com/sameh-farouk/cometbft-custom-app.git
cd cometbft-custom-app

# Build with default database (Badger, pebble)
make build-linux

# start the 4 validators local network
DB_TYPE=badger make start-localnet
# or
DB_TYPE=pebble make start-localnet

# apply latency to simulate a real-world network environment
make apply-latency

# when you done
make stop-localnet

# clean up
sudo make clean 
```

## Using Different Database Backends

You can specify which database to use with the `DB_TYPE` environment variable, but TigerBeetle DB requires a specific build tag, so you need to build the binary with the `tigerbeetle` build tag.

```bash
# Build with TigerBeetle DB support
make build-tigerbeetle
DB_TYPE=tigerbeetle make start-localnet
```

Note: TigerBeetle DB support is currently WIP.

## Example Usage

There are 4 accounts pre-created: 1, 2, 3, 4
private keys are:

- "1": "23e980b97c67af9b94319b6672049fbd2f9992eaf6a567a2b5a66286e527e8e9c8af5ee74756bb934c9c3f93a3ffa4125c93d8a76619a1834f4511334d83d45f"
- "2": "11a2070b5bf25002c43d238117840fb97492266d3e0fb7637b069d5569b5d8283382d764d3e30ce4c3aab066335a558e8f632d2aaf161e6aa5615c57176cfbca"
- "3": "376384a0c4d4ef4e95bf980acea6ec6d7b8bbdaa06d91ca68383d018e885dda204c01c7d4f6c784504fce83f97968145e8aa6ca461ec19f3a685466152f17644"
- "4": "7403b7706deba2d8d036d00c5e1e087542fff733b1b3f1b776bf2fa64bcd5d98d06a22ce4b7a59ceac3a898504901f41e27491ed3cc90e8ee46ac43e9305d61a"

```bash
# Send a transaction tx_id=1, sender_id=1, receiver_id=2, amount=50, signature=<SIGNATURE>
curl -s 'localhost:26657/broadcast_tx_commit?tx="1=1=2=50=<SIGNATURE>"'

# Query account balances
curl -s 'localhost:26657/abci_query?data="1"'
echo OTk5OTk5OTUw | base64 --decode  # 999999950

curl -s 'localhost:26657/abci_query?data="2"'
echo MTAwMDAwMDA1MA== | base64 --decode  # 1000000050

# you can send batched transactions as well separated by `:`
curl -s 'localhost:26657/broadcast_tx_commit?tx="2=1=2=50=<SIGNATURE>:3=1=2=50=<SIGNATURE>:4=1=2=50=<SIGNATURE>"'
```

## Database Configuration

Each database can be configured with additional options:

- **Badger DB**: Use `-db-path` to specify the database directory
- **Pebble DB**: Use `-db-path` to specify the database directory
- **TigerBeetle DB**: 
  - Use `-tb-addresses` to specify TigerBeetle server addresses (comma-separated)
  - Use `-tb-cluster-id` to specify the TigerBeetle cluster ID
