# cryptocoins
Research all coins hash, address and signature algorithm.

| No | Name | Website | Explorer | Source code | Signature | Sig2 | Hash | Address | Address format | bip44 | testnet rpc | example | reference |
|-----|------------------|------------------------------|--------------------------------------|-------------------------------------|---------------------|--------------------|----------------|----------------------------------------------------------------------------------------------------------|----------------|-------|--------------------------------------------------------|---------|-----------------------------------------------------------------------------------------------------------------------------------|
| 1 | Bitcoin | https://bitcoin.org/ | https://btc.com/ | https://github.com/bitcoin/bitcoin | ECDSA/secp256k1 |  | SHA256 | 137H4GbcDz3e3DS9ADDbee4wa1GHHdcEnW | base58 |  |  |  |  |
| 2 | Ethereum | https://www.ethereum.org/ | https://etherscan.io/ | https://github.com/ethereum/ | ECDSA/secp256k1 |  | SHA3-keccak256 | 0x0bDcBCbB9B0aCA3EAEE7a94A4fb5FB0f16681e2f | hex |  |  |  |  |
| 3 | XRP | https://ripple.com/ | https://xrpcharts.ripple.com/#/ | https://github.com/ripple | ECDSA/secp256k1 | ED25519/Curve25519 | SHA256 | rUnpzXPagSPLE2CbTsxa5Ey2xGy62PmDem | base58 |  | https://developers.ripple.com/xrp-test-net-faucet.html |  |  |
| 4 | Bitcoin Cash | https://www.bitcoincash.org/ | https://bch.btc.com/ | https://github.com/bitcoincashorg/ | ECDSA/secp256k1 |  | SHA256 | 1MSh2kijYoZr4cHAVQiNSrYDzZcMhjUYk1 | base58 |  |  |  |  |
| 5 | Stellar | https://www.stellar.org/ | https://stellarchain.io/ | https://github.com/stellar/ | ED25519/Curve25519 |  | SHA256 | GAI3GJ2Q3B35AOZJ36C4ANE3HSS4NK7WI6DNO4ZSHRAX6NG7BMX6VJER |  |  |  |  |  |
| 6 | EOS | https://eos.io/ | https://eospark.com/ | https://github.com/eosio | ECDSA/secp256k1 |  | SHA256 | EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV |  |  |  |  |  |
| 7 | Litecoin | https://litecoin.org/ | https://live.blockcypher.com/ltc/ | https://github.com/litecoin-project | ECDSA/secp256k1 |  | SHA256 | LS3L3ThLEyuMWHkDts9BgXseu1byLgbpDw |  |  |  |  |  |
| 8 | Cardano | https://www.cardano.org | https://cardanoexplorer.com/ | https://github.com/input-output-hk | ED25519/Curve25519 |  | Blake2b-224 | DdzFFzCqrhshvNLoWJi1uYoRoHYPDJrv1gzd4CGgKSYHWiqbK8RDiWwnSp9iTwpFENBeRbtJ5dLQURouHntANNKQLAfUsPBRWT3pWj5r | base58 |  |  |  |  |
| 9 | Monero | https://getmonero.org/ | https://moneroblocks.info/ | https://github.com/monero-project |  |  |  | 49Jt4tzbvZ5PyEMub6tNDGKP4zxogN9t1VACVWgTEcMwhtCGjxrDyt5XCDHG6XpA2U1uWsnsyKYdrL25Vp6y2pou2bdboCZ |  |  |  |  |  |
| 10 | Tether | https://tether.to/ | https://www.omniexplorer.info/ | https://github.com/OmniLayer/ | ECDSA/secp256k1 |  | SHA256 | 18KmBuZVAK9qMq38gm5DwFkZ7asvuhGyFT |  |  |  |  |  |
| 11 | TRON | https://tron.network/ | https://tronscan.org | https://github.com/tronprotocol | ECDSA/secp256k1 |  | SHA3-keccak256 | TMYcx6eoRXnePKT1jVn25ZNeMNJ6828HWk |  |  |  |  | https://github.com/tronprotocol/Documentation/blob/master/English_Documentation/Procedures_of_transaction_signature_generation.md |
| 12 | IOTA | https://www.iota.org/ | https://thetangle.org/live | https://github.com/iotaledger | Winternitz one-time |  | keccak384 | JJYINSNHNLDVI9P9HITKMSKJKMXTKDVIULWRCFCBNPEKMYBD9DLSKHMNIYZBSBQFLIYRBSC9ZXDJAESMCVTYQPQDRY |  |  |  |  | https://iota.readme.io/docs/seeds-private-keys-and-accounts |
| 13 | Dash | https://www.dash.org/ | https://explorer.dash.org/chain/Dash | https://github.com/dashpay/ | ECDSA/secp256k1 |  | SHA256 | Xubdr2uaECHfyhVBnAoPCX4dxv14yR3W4d |  |  |  |  |  |
| 14 | NEO | https://neo.org/ | http://antcha.in/ | https://github.com/neo-project | ECDSA/secp256r1 |  |  | AKCbHhCf3Sq9qeCm8n2nmhGdgMuTDrEhmK |  |  |  |  | https://github.com/neo-project/neo/blob/master/neo/Wallets/KeyPair.cs |
| 15 | Ethereum Classic | https://ethereumclassic.org/ | https://etherhub.io/home | https://github.com/ethereumproject | ECDSA/secp256k1 |  | SHA3-keccak256 | 0x18489e2a517b22348f20343b7386b6a81d5261c4 |  |  |  |  |  |
| 16 | NEM | https://nem.io/ | http://explorer.nemchina.com | https://github.com/NemProject | ED25519/Curve25519 |  |  | NDNBRZZ3VZGZ626NFVL357APEYACNL6NMTRKTF5W |  |  |  |  | https://nem.io/wp-content/themes/nem/files/NEM_techRef.pdf |
| 17 | Tezos | https://tezos.com/ | https://tezos.id | https://gitlab.com/tezos/tezos |  |  |  | tz1bZ8vsMAXmaWEV7FRnyhcuUs2fYMaQ6Hkk |  |  |  |  |  |
| 18 | Zcash | https://z.cash/ | https://explorer.zcha.in/ | https://github.com/zcash/ | ECDSA/secp256k1 | ED25519/Curve25519 |  | t1aZ2DGuiokCxHVfb4cGQqXghxy9hUpE6xQ |  |  |  |  |  |
| 19 | VeChain | https://www.vechain.org/ | https://explore.veforge.com/ | https://github.com/vechain/ | ECDSA/secp256k1 |  | SHA3-keccak256 | 0xdde1C7AD4Cca5672a5c6DB767B7ed79794bF7ca8 |  |  |  |  |  |
| 20 | Bitcoin Gold |  |  |  |  |  |  |  |  |  |  |  |  |
| 21 | Decred |  |  |  |  |  |  |  |  |  |  |  |  |
| 22 | Dogecoin |  |  |  |  |  |  |  |  |  |  |  |  |
| 23 | Qtum |  |  |  |  |  |  |  |  |  |  |  |  |
| 24 | Ontology |  |  |  |  |  |  |  |  |  |  |  |  |
| 25 | Lisk |  |  |  |  |  |  |  |  |  |  |  |  |
| 26 | Bitcoin Diamond |  |  |  |  |  |  |  |  |  |  |  |  |
| 27 | ICON |  |  |  |  |  |  |  |  |  |  |  |  |
| 28 | BitShares |  |  |  |  |  |  |  |  |  |  |  |  |
| 29 | Bytecoin |  |  |  |  |  |  |  |  |  |  |  |  |
| 30 | Nano |  |  |  |  |  |  |  |  |  |  |  |  |
| 31 | DigiByte |  |  |  |  |  |  |  |  |  |  |  |  |
| 32 | Siacoin |  |  |  |  |  |  |  |  |  |  |  |  |
| 33 | Steem |  |  |  |  |  |  |  |  |  |  |  |  |
| 34 | Verge |  |  |  |  |  |  |  |  |  |  |  |  |
| 35 | Bytom |  |  |  |  |  |  |  |  |  |  |  |  |
| 36 | Waves |  |  |  |  |  |  |  |  |  |  |  |  |
| 37 | Metaverse ETP |  |  |  |  |  |  |  |  |  |  |  |  |
| 38 | Stratis |  |  |  |  |  |  |  |  |  |  |  |  |
| 39 | Komodo |  |  |  |  |  |  |  |  |  |  |  |  |
| 40 | Electroneum |  |  |  |  |  |  |  |  |  |  |  |  |
| 41 | Cryptonex |  |  |  |  |  |  |  |  |  |  |  |  |
| 42 | Ardor |  |  |  |  |  |  |  |  |  |  |  |  |
| 43 | Wanchain |  |  |  |  |  |  |  |  |  |  |  |  |
| 44 | MOAC |  |  |  |  |  |  |  |  |  |  |  |  |
| 45 | Ravencoin |  |  |  |  |  |  |  |  |  |  |  |  |
| 46 | GXChain |  |  |  |  |  |  |  |  |  |  |  |  |
| 47 | Huobi Token |  |  |  |  |  |  |  |  |  |  |  |  |
| 48 | PIVX |  |  |  |  |  |  |  |  |  |  |  |  |
| 49 | MonaCoin |  |  |  |  |  |  |  |  |  |  |  |  |
| 50 | Ark |  |  |  |  |  |  |  |  |  |  |  |  |
| 51 | Horizen |  |  |  |  |  |  |  |  |  |  |  |  |
| 52 | HyperCash |  |  |  |  |  |  |  |  |  |  |  |  |
| 53 | ReddCoin |  |  |  |  |  |  |  |  |  |  |  |  |
| 54 | Nxt |  |  |  |  |  |  |  |  |  |  |  |  |
| 55 | Elastos |  |  |  |  |  |  |  |  |  |  |  |  |
| 56 | Zcoin |  |  |  |  |  |  |  |  |  |  |  |  |
| 57 | Nebulas |  |  |  |  |  |  |  |  |  |  |  |  |
| 58 | GoChain |  |  |  |  |  |  |  |  |  |  |  |  |
| 59 | BOScoin |  |  |  |  |  |  |  |  |  |  |  |  |
| 60 | Syscoin |  |  |  |  |  |  |  |  |  |  |  |  |
| 61 | Bitcoin Private |  |  |  |  |  |  |  |  |  |  |  |  |
| 62 | Factom |  |  |  |  |  |  |  |  |  |  |  |  |
| 63 | BridgeCoin |  |  |  |  |  |  |  |  |  |  |  |  |
| 64 | WaykiChain |  |  |  |  |  |  |  |  |  |  |  |  |
| 65 | Emercoin |  |  |  |  |  |  |  |  |  |  |  |  |
| 66 | Peercoin |  |  |  |  |  |  |  |  |  |  |  |  |
| 67 | Nexus |  |  |  |  |  |  |  |  |  |  |  |  |
| 68 | Groestlcoin |  |  |  |  |  |  |  |  |  |  |  |  |
| 69 | Byteball Bytes |  |  |  |  |  |  |  |  |  |  |  |  |
| 70 | Neblio |  |  |  |  |  |  |  |  |  |  |  |  |
| 71 | Vertcoin |  |  |  |  |  |  |  |  |  |  |  |  |
| 72 | Skycoin |  |  |  |  |  |  |  |  |  |  |  |  |
| 73 | TokenPay |  |  |  |  |  |  |  |  |  |  |  |  |
| 74 | Einsteinium |  |  |  |  |  |  |  |  |  |  |  |  |
| 75 | Ubiq |  |  |  |  |  |  |  |  |  |  |  |  |
| 76 | Blocknet |  |  |  |  |  |  |  |  |  |  |  |  |
| 77 | SmartCash |  |  |  |  |  |  |  |  |  |  |  |  |
| 78 | Apollo Currency |  |  |  |  |  |  |  |  |  |  |  |  |
| 79 | SaluS |  |  |  |  |  |  |  |  |  |  |  |  |
| 80 | NavCoin |  |  |  |  |  |  |  |  |  |  |  |  |
| 81 | Haven Protocol |  |  |  |  |  |  |  |  |  |  |  |  |
| 82 | Novacoin |  |  |  |  |  |  |  |  |  |  |  |  |
| 83 | Nasdacoin |  |  |  |  |  |  |  |  |  |  |  |  |
| 84 | DigitalNote |  |  |  |  |  |  |  |  |  |  |  |  |
| 85 | POA Network |  |  |  |  |  |  |  |  |  |  |  |  |
| 86 | Achain |  |  |  |  |  |  |  |  |  |  |  |  |
| 87 | BitBay |  |  |  |  |  |  |  |  |  |  |  |  |
| 88 | Viacoin |  |  |  |  |  |  |  |  |  |  |  |  |
| 89 | Vitae |  |  |  |  |  |  |  |  |  |  |  |  |
| 90 | WhiteCoin |  |  |  |  |  |  |  |  |  |  |  |  |
| 91 | Particl |  |  |  |  |  |  |  |  |  |  |  |  |
| 92 | Loki |  |  |  |  |  |  |  |  |  |  |  |  |
| 93 | NIX |  |  |  |  |  |  |  |  |  |  |  |  |
| 94 | Burst |  |  |  |  |  |  |  |  |  |  |  |  |
| 95 | Pascal Coin |  |  |  |  |  |  |  |  |  |  |  |  |
| 96 | E-Dinar Coin |  |  |  |  |  |  |  |  |  |  |  |  |
| 97 | Nexty |  |  |  |  |  |  |  |  |  |  |  |  |
| 98 | Steem Dollars |  |  |  |  |  |  |  |  |  |  |  |  |
| 99 | ALQO |  |  |  |  |  |  |  |  |  |  |  |  |
| 100 | CloakCoin |  |  |  |  |  |  |  |  |  |  |  |  |
| 101 | Energi |  |  |  |  |  |  |  |  |  |  |  |  |

