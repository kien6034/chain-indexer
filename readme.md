# BTC - EVM Integration

- [x] BTC indexer by wallet transactions
- [ ] BTC signing client: https://github.com/btcsuite/btcwallet
- [ ] EVM indexer / signing client

## Chain-indexer

Indexer for btc

## Methods

### Get Address Transactions

```go
type Transaction struct {
  BaseAddress  string     `json:"base_address"`
	TxId         string     `json:"txid"`
	Sender       []Sender   `json:"sender"` // TODO: sender could be multiple addresses
	Receivers    []Receiver `json:"receivers"`
	TotalSpend   int64      `json:"total_spend"`
	TotalReceive int64      `json:"total_receive"`
	TxType       TxType     `json:"tx_type"`
}

func (c* Client) GetAddressTransactions(address string) ([]Transaction, error)
```

Example output:

```
 === Tx: 6d5cac4fe973e69628115af52608487238700f6c28cd0f992591da027c6ec4ff ===
Base address: tb1qjfaa5vvxt9m4sp9kqkcpzypkzydz2vcywqx9tm->> Outgoing: Sent 3000 sats from tb1qjfaa5vvxt9m4sp9kqkcpzypkzydz2vcywqx9tm

Receiver: tb1p6g4krnjtkrdmxa9mky8n05r9xntv5qnf3gl84adrr5wpr6w7ga7q4tggyz || Amount 400
Receiver: tb1psc68fp24y54zfmw5pfuj0uwsqp8nwpw7cjpprsx98q77298mp99s2ajttw || Amount 2435

 === Tx: 410267bc9fc0bd77313bab0fdc2ee7ca48b0c5e6235ed3c9010c8e5b16804113 ===
Base address: tb1qjfaa5vvxt9m4sp9kqkcpzypkzydz2vcywqx9tm->> Outgoing: Sent 800000 sats from tb1qjfaa5vvxt9m4sp9kqkcpzypkzydz2vcywqx9tm

Receiver: tb1p6g4krnjtkrdmxa9mky8n05r9xntv5qnf3gl84adrr5wpr6w7ga7q4tggyz || Amount 400
Receiver: tb1pgpytdxfqukzmcyh3q0m0v3t6kaqccq4ttspjz2qdp72k9np6h9ps7d09t7 || Amount 799435

 === Tx: dacbc695d656d027fbf16ad19303a8c7080ea22336758a2ad47e274e4e9f1b1f ===
Base address: tb1qjfaa5vvxt9m4sp9kqkcpzypkzydz2vcywqx9tm<<- Incoming: Received 800000 sats to tb1qjfaa5vvxt9m4sp9kqkcpzypkzydz2vcywqx9tm

 Sender: tb1qhh7tfvdmfuyu2py0xjkcs6fpyz4t9ha43letz4 || Amount 800685

 === Tx: 84f6a2088b8cb9d50cd487a22ed2d6f18effd8335a38574d63062cce3f060864 ===
Base address: tb1qjfaa5vvxt9m4sp9kqkcpzypkzydz2vcywqx9tm<<- Incoming: Received 3000 sats to tb1qjfaa5vvxt9m4sp9kqkcpzypkzydz2vcywqx9tm

 Sender: tb1qy4q88efmcxjm6u3z7c8l8thydhujngqjumll00 || Amount 965456488

 === Tx: cea31db822dae0b985fbc0d84f9eca75db259a1a5291f52e34f0ec056bd10cdf ===
Base address: tb1qjfaa5vvxt9m4sp9kqkcpzypkzydz2vcywqx9tm->> Outgoing: Sent 1000 sats from tb1qjfaa5vvxt9m4sp9kqkcpzypkzydz2vcywqx9tm

Receiver: tb1p6g4krnjtkrdmxa9mky8n05r9xntv5qnf3gl84adrr5wpr6w7ga7q4tggyz || Amount 1000
Receiver: tb1pm3q6haw23t3mgs0lntn7pardz7ljmsn3slvzvumxd6gwzexg8wlqhq4wk5 || Amount 613

 === Tx: aa0f65543dd78400ef7ced418cc08e7aaf4604aab8f77889f12fbd1d62c90015 ===
Base address: tb1qjfaa5vvxt9m4sp9kqkcpzypkzydz2vcywqx9tm->> Outgoing: Sent 2000 sats from tb1qjfaa5vvxt9m4sp9kqkcpzypkzydz2vcywqx9tm

Receiver: tb1p6g4krnjtkrdmxa9mky8n05r9xntv5qnf3gl84adrr5wpr6w7ga7q4tggyz || Amount 1000
Receiver: tb1ph7r7rnn9d0kc65hv9nq64xnknc4g32saf78x8elm3axmw6ef4qpqcsgrg3 || Amount 835

 === Tx: de9d9ac3b8c38159b1f1e4d1aea40c5ec4db2da2d2bdc1749ff4f9bcd05bd09c ===
Base address: tb1qjfaa5vvxt9m4sp9kqkcpzypkzydz2vcywqx9tm<<- Incoming: Received 1000 sats to tb1qjfaa5vvxt9m4sp9kqkcpzypkzydz2vcywqx9tm

 Sender: tb1p5mz3tvruxm58dkss6yvrs3zygng2k6je9zz4f3t7uhvce2z833gseq0utg || Amount 960364553

 === Tx: ff1ce773135ca1a520f4a56f50d159923c574d4a888ec13064527b7f5e2a729a ===
Base address: tb1qjfaa5vvxt9m4sp9kqkcpzypkzydz2vcywqx9tm<<- Incoming: Received 2000 sats to tb1qjfaa5vvxt9m4sp9kqkcpzypkzydz2vcywqx9tm

 Sender: tb1pglnp34zmrtchtqdmvrh9w6esp3j3xqgvmtct2hgyz6u235dp9gyqgtjt78 || Amount 942528753
```

## Todos

- [ ] Add logging
- [ ] Save tx to db / cache
