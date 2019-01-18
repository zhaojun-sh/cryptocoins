# golang version 
### clone and build demo 
#### 1. ensure dependences:  
`github.com/btcsuite/btcd`  
`github.com/btcsuite/btcutil`  
`github.com/btcsuite/btcwallet`  
`github.com/ethereum/go-ethereum`  
`github.com/eoscanada/eos-go`  
`github.com/rubblelabs/ripple`  
#### 2. clone this project 
`$ git clone https://github.com/gaozhengxin/cryptocoins.git`  
#### 3. build demo program 
`$ cd ./cryptocoins/src/go`  
`$ go build ./demo/main.go`  
#### 4. run demo program 
`$ ./main`  

  
  
  
### add support for new cryptocurrency 
#### 1. 
Build a package in `src/go`, and write your code in it. You are supposed to define a struct that implements the interface  TransactionHandler. You can find the interface definition in `src/go/api.go`. Configuration constants such as the urls of gateways should be defined in package `src/go/config`.
#### 2. 
Register new transaction handler in `src/go/api.go`. Insert the constructor of your transaction handler is the switch-case statement.
