# Quick start

```bash
git clone https://github.com/sameh-farouk/cometbft-custom-app.git
cd cometbft-custom-app

make build-linux
make start-localnet
make apply-latency
```

```bash
> curl -s 'localhost:26657/broadcast_tx_commit?tx="1=2=50=93038a25e2806ea16771281f734b8dc28600082ef05f59e39a6765c5f73dd115991dd91e5c4af11e7d20b5170519eee281d873dcb625730e95d8e0fb3ff23902"'
{"jsonrpc":"2.0","id":-1,"result":{"check_tx":{"code":0,"data":null,"log":"","info":"","gas_wanted":"0","gas_used":"0","events":[],"codespace":""},"tx_result":{"code":0,"data":null,"log":"","info":"","gas_wanted":"0","gas_used":"0","events":[{"type":"app","attributes":[{"key":"src","value":"1","index":true},{"key":"dst","value":"2","index":true},{"key":"amount","value":"50","index":true}
> curl -s 'localhost:26657/abci_query?data="1"'
{"jsonrpc":"2.0","id":-1,"result":{"response":{"code":0,"log":"exists","info":"","index":"0","key":"MQ==","value":"OTk5OTk5OTUw","proofOps":null,"height":"0","codespace":""}}}
> echo OTk5OTk5OTUw | base64 --decode
999999950
> curl -s 'localhost:26657/abci_query?data="2"'
{"jsonrpc":"2.0","id":-1,"result":{"response":{"code":0,"log":"exists","info":"","index":"0","key":"Mg==","value":"MTAwMDAwMDA1MA==","proofOps":null,"height":"0","codespace":""}}}
> echo MTAwMDAwMDA1MA== | base64 --decode
1000000050
```
