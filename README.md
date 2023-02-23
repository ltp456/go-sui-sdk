# go-sui-sdk
sui rest api implementation of the go language version, generate address, sign transaction and others

[SUI Node Api Doc](https://docs.sui.io/sui-jsonrpc#sui_batchTransaction)

## example


use example

    package main

    
    func main() {
    
        endpoint := "https://fullnode.devnet.sui.io:443"
        //endpoint := "http://127.0.0.1:9000"
        client, err := NewSuiClient(endpoint)
        if err != nil {
            panic(err)
        }

        number, err := client.GetTotalTransactionNumber()
        if err != nil {
            panic(err)
        }
        fmt.Println(number)

    }

## Todo