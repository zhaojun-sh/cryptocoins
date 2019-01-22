# golang version 
### clone and build demo   
#### 1. download and build 
`$ go get -u -d github.com/gaozhengxin/cryptocoins/src/go`  
#### 2. compile the demo program 
`$ cd $GOPATH/src/github.com/gaozhengxin/cryptocoins/src/go`  
`$ go build ./demo/main.go`  
#### 3. run demo program 
`$ ./main`  

  
  
  
### add support for new cryptocurrency 
#### 1. 
Build a package in `src/go`, and write your code in it. You are supposed to define a struct that implements the interface  TransactionHandler. You can find the interface definition in `src/go/api.go`. Configuration constants such as the urls of gateways should be defined in package `src/go/config`.
#### 2. 
Register new transaction handler in `src/go/api.go`. Insert the constructor of your transaction handler is the switch-case statement.
