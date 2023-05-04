# go-sui-sdk

#### example

    //endpoint := "https://fullnode.devnet.sui.io:443"
	endpoint := "http://127.0.0.1:9000"
	client, err = NewSuiClient(endpoint)
	if err != nil {
		panic(err)
	}