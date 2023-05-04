package go_sui_sdk

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"go-sui-sdk/crypto"
	"go-sui-sdk/types"
	"math/big"
	"testing"
)

var client *SuiClient
var err error

var endpoint = "https://fullnode.mainnet.sui.io:443"

//var endpoint = "https://fullnode.testnet.sui.io:443"

func init() {
	//endpoint := "https://fullnode.devnet.sui.io:443"
	//endpoint := "http://127.0.0.1:9000"
	client, err = NewSuiClient(endpoint)
	if err != nil {
		panic(err)
	}
}

func TestSuiClient_ScanBlock(t *testing.T) {
	latestCheckpointSequenceNumber, err := client.GetLatestCheckpointSequenceNumber()
	if err != nil {
		panic(err)
	}
	for index := latestCheckpointSequenceNumber; ; index = index + 1 {
		fmt.Println(index)
		transactions, _, err := client.Transactions(index, 1)
		if err != nil {
			panic(err)
		}
		for _, tx := range transactions {
			fmt.Println(tx)
		}
	}

}

func TestSuiClient_GetCoins(t *testing.T) {
	coins, err := client.GetCoins(types.SuiCoinType, "0xa2909dd355aeaffd520847a35da1ef529c4cebd348bb4d69f00f32315118106e", "")
	if err != nil {
		panic(err)
	}
	fmt.Println(coins)
}

func TestSuiClient_GetEvents(t *testing.T) {
	events, err := client.GetEvents("CfMK9eW76jwGD3iQHoGiH9Gr1iwp8k2ntuFLDRT9Hbs1")
	if err != nil {
		panic(err)
	}
	fmt.Println(events)
}

func TestSuiClient_GetReferenceGasPrice(t *testing.T) {
	price, err := client.GetReferenceGasPrice()
	if err != nil {
		panic(err)
	}
	fmt.Println(price)
}

func TestSuiClient_GetObject(t *testing.T) {
	object, err := client.GetObject("0x5994bda6a98e5b8f29717bb066cf2b309344c1aa6cf247ed8ab90e244857394b")
	if err != nil {
		panic(err)
	}
	fmt.Println(object)
}

func TestSuiClient_GetOwnedObjects(t *testing.T) {
	objects, err := client.GetOwnedObjects("", "0xa2909dd355aeaffd520847a35da1ef529c4cebd348bb4d69f00f32315118106e")
	if err != nil {
		panic(err)
	}
	for _, item := range objects.Data {
		fmt.Println(item.Data.ObjectID)
	}
}

func TestSuiClient_GetBalance(t *testing.T) {
	balance, err := client.GetBalance(types.SuiCoinType, "0xc0ee0c49b3be532975fdb2c02a3ae8dea70b58f53879f70ba256974627e23ee3")
	if err != nil {
		panic(err)
	}
	fmt.Println(balance)
}

func TestSuiClient_GetAllBalances(t *testing.T) {
	allBalances, err := client.GetAllBalances("0xc0ee0c49b3be532975fdb2c02a3ae8dea70b58f53879f70ba256974627e23ee3")
	if err != nil {
		panic(err)
	}
	fmt.Println(allBalances)
}

func TestSuiClient_GetTx(t *testing.T) {
	transactionBlock, err := client.GetTransactionBlock("HrEB6m8Qv2mrHqs6G82r8wsXacGTV4cmppJCmFi4PTmc")
	if err != nil {
		panic(err)
	}
	fmt.Println(transactionBlock.Status())
	fmt.Println(transactionBlock.GetGasUsed())
	txes, err := transactionBlock.Parse()
	if err != nil {
		panic(err)
	}
	for _, tx := range txes {
		fmt.Println(tx)
	}
}

func TestSuiClient_GetTransactionBlock(t *testing.T) {
	digestList := []string{
		"6RdJ69bh7LKWUa8upf1gRxjCApDhgvJ9S2pXDmdUHPeW",
		"3QixuaFW8BJWK1Sg2H7z69p3yF9tp28k7ysAd8Wg6UNw",
		"5KxLweSkBg3uavjjyMNuTdwEP42Kqhyjm6KnCciKpWs6",
		"FgnntShXMos58nGvjsR1ZtjgXHunR1PmGbhBJ1HP9J7U",
		"Cp7XhffS9uoenKRnavGuXUiF4zrhqMNSxrMmD8YonuG5",
		"65DUE8n3uMH8RHeb7RfHFFsVatXV4Tcmrsu87aH11rju",
	}
	for _, digest := range digestList {
		transactionBlock, err := client.GetTransactionBlock(digest)
		if err != nil {
			panic(err)
		}
		txList, err := transactionBlock.Parse()
		if err != nil {
			panic(err)
		}
		for _, tx := range txList {
			fmt.Println(tx)
		}
	}

}

func TestSuiClient_TransactionsByCheckpoint(t *testing.T) {
	transactions, _, err := client.Transactions(2038639, 1)
	if err != nil {
		panic(err)
	}
	for _, item := range transactions {
		fmt.Println(item)
	}
}

func TestSuiClient_GetCheckpoints(t *testing.T) {
	checkpoints, err := client.GetCheckpoints(1906812, 10)
	if err != nil {
		panic(err)
	}
	for _, item := range checkpoints.Data {
		fmt.Println(item)
	}
}

func TestSuiClient_GetCheckpoint(t *testing.T) {
	checkpoint, err := client.GetCheckpoint(1906812)
	if err != nil {
		panic(err)
	}
	fmt.Println(checkpoint)
}

func TestSuiClient_GetLatestCheckpointSequenceNumber(t *testing.T) {
	getLatestCheckpointSequenceNumber, err := client.GetLatestCheckpointSequenceNumber()
	if err != nil {
		panic(err)
	}
	fmt.Println(getLatestCheckpointSequenceNumber)
}

func TestSuiClient_GetTotalTransactionNumber(t *testing.T) {
	number, err := client.GetTotalTransactionBlocks()
	if err != nil {
		panic(err)
	}
	fmt.Println(number)
}

func TestSuiClient_Transfer(t *testing.T) {
	/*
		seed:dcee1121d8e62b67f9a5e8a9ebb49fe1d3b41aaf7bee77984ac4648e632cddc0
		privateKey:dcee1121d8e62b67f9a5e8a9ebb49fe1d3b41aaf7bee77984ac4648e632cddc0fd7d7902f8e171d460fac5cfadb0b92f299e5bee16f5a1d454c2ed190d0bbb82
		publicKey:fd7d7902f8e171d460fac5cfadb0b92f299e5bee16f5a1d454c2ed190d0bbb82
		addess:0xc0ee0c49b3be532975fdb2c02a3ae8dea70b58f53879f70ba256974627e23ee3
	*/

	sender := "0xc0ee0c49b3be532975fdb2c02a3ae8dea70b58f53879f70ba256974627e23ee3"
	recipent := "0xa2909dd355aeaffd520847a35da1ef529c4cebd348bb4d69f00f32315118106e"
	amount := "1110000000"
	gasBudegt := "30000000"
	allObjectIds, err := client.GetAllCoinObjectIds(types.SuiCoinType, sender)
	if err != nil {
		panic(err)
	}
	transactionBlock, err := client.DryRunTransactionBlock(types.SuiCoinType, true, sender, allObjectIds, recipent, amount, gasBudegt)
	if err != nil {
		panic(err)
	}
	fmt.Println(transactionBlock.Status())
	gasUsed, err := transactionBlock.GetGasUsed()
	if err != nil {
		panic(err)
	}
	fmt.Println(gasUsed)

	return
	seed := "dcee1121d8e62b67f9a5e8a9ebb49fe1d3b41aaf7bee77984ac4648e632cddc0"
	seedBytes, err := hex.DecodeString(seed)
	if err != nil {
		panic(err)
	}
	amountBig, ok := big.NewInt(0).SetString(amount, 10)
	if !ok {
		panic("ddsdfs")
	}
	balanceBig, err := client.Balance(types.SuiCoinType, sender)
	if err != nil {
		panic(err)
	}
	amountBig = big.NewInt(0).Sub(balanceBig, gasUsed)

	tx, err := client.Transfer(types.SuiCoinType, seedBytes, sender, allObjectIds, []string{recipent}, []string{amountBig.String()}, gasUsed.String())
	if err != nil {
		panic(err)
	}
	fmt.Println(tx.Digest)
}

func TestSuiClient_DevInspectTransactionBlock(t *testing.T) {
	sender := "0xc0ee0c49b3be532975fdb2c02a3ae8dea70b58f53879f70ba256974627e23ee3"
	recipent := "0xa2909dd355aeaffd520847a35da1ef529c4cebd348bb4d69f00f32315118106e"
	amount := "2300000"
	gasBudget := "0"
	allObjectIds, err := client.GetAllCoinObjectIds(types.SuiCoinType, sender)
	if err != nil {
		panic(err)
	}
	result, err := client.DevInspectTransactionBlock(types.SuiCoinType, false, sender, allObjectIds, recipent, amount, gasBudget)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.GetGasUsed())
	fmt.Println(result.Status())
}

func TestSuiClient_DryRunTransactionBlock(t *testing.T) {
	sender := "0xc0ee0c49b3be532975fdb2c02a3ae8dea70b58f53879f70ba256974627e23ee3"
	recipent := "0xa2909dd355aeaffd520847a35da1ef529c4cebd348bb4d69f00f32315118106e"
	amount := "1100000000"
	gasBudget := "1100000000"
	allObjectIds, err := client.GetAllCoinObjectIds(types.SuiCoinType, sender)
	if err != nil {
		panic(err)
	}
	result, err := client.DryRunTransactionBlock(types.SuiCoinType, true, sender, allObjectIds, recipent, amount, gasBudget)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.GetGasUsed())
	fmt.Println(result.Status())
}

func TestSuiClient_Pay(t *testing.T) {

	seed := "dcee1121d8e62b67f9a5e8a9ebb49fe1d3b41aaf7bee77984ac4648e632cddc0"
	seedBytes, err := hex.DecodeString(seed)
	if err != nil {
		panic(err)
	}
	sender := "0xc0ee0c49b3be532975fdb2c02a3ae8dea70b58f53879f70ba256974627e23ee3"
	recipent := []string{"0xa2909dd355aeaffd520847a35da1ef529c4cebd348bb4d69f00f32315118106e", "0x3da5a91eab1be9ef35e5b6fe65ed9328e08e23cdf9cc20b7131ff0d095b97e9f"}
	objectIds, err := client.GetAllCoinObjectIds(types.SuiCoinType, sender)
	if err != nil {
		panic(err)
	}
	gasObjectId := objectIds[0]
	objectIds = objectIds[1:]
	amount := []string{"1100000", "1200000"}
	gasBudget := "4000000"
	tx, err := client.Pay(seedBytes, sender, gasObjectId, objectIds, recipent, amount, gasBudget)
	if err != nil {
		panic(err)
	}
	fmt.Println(tx.Digest)
	//6RdJ69bh7LKWUa8upf1gRxjCApDhgvJ9S2pXDmdUHPeW
	//65DUE8n3uMH8RHeb7RfHFFsVatXV4Tcmrsu87aH11rju

}

func TestSuiClient_PaySui(t *testing.T) {
	seed := "dcee1121d8e62b67f9a5e8a9ebb49fe1d3b41aaf7bee77984ac4648e632cddc0"
	seedBytes, err := hex.DecodeString(seed)
	if err != nil {
		panic(err)
	}
	sender := "0xc0ee0c49b3be532975fdb2c02a3ae8dea70b58f53879f70ba256974627e23ee3"
	recipent := "0xa2909dd355aeaffd520847a35da1ef529c4cebd348bb4d69f00f32315118106e"
	objectIds, err := client.GetAllCoinObjectIds(types.SuiCoinType, sender)
	if err != nil {
		panic(err)
	}
	amount := "1100000"
	gasBudget := "40000000"

	tx, err := client.PaySui(seedBytes, sender, objectIds, []string{recipent}, []string{amount}, gasBudget)
	if err != nil {
		panic(err)
	}
	fmt.Println(tx.Digest)
	//3QixuaFW8BJWK1Sg2H7z69p3yF9tp28k7ysAd8Wg6UNw
}

func TestSuiClient_PayAllSui(t *testing.T) {
	seed := "dcee1121d8e62b67f9a5e8a9ebb49fe1d3b41aaf7bee77984ac4648e632cddc0"
	seedBytes, err := hex.DecodeString(seed)
	if err != nil {
		panic(err)
	}
	sender := "0xc0ee0c49b3be532975fdb2c02a3ae8dea70b58f53879f70ba256974627e23ee3"
	recipent := "0xa2909dd355aeaffd520847a35da1ef529c4cebd348bb4d69f00f32315118106e"
	objectIds, err := client.GetAllCoinObjectIds(types.SuiCoinType, sender)
	if err != nil {
		panic(err)
	}
	gasBudget := "3000000"
	tx, err := client.PayAllSui(seedBytes, sender, recipent, objectIds, gasBudget)
	if err != nil {
		panic(err)
	}
	fmt.Println(tx.Digest)
	//5KxLweSkBg3uavjjyMNuTdwEP42Kqhyjm6KnCciKpWs6
}

func TestSuiClient_TransferSui(t *testing.T) {

	seed := "dcee1121d8e62b67f9a5e8a9ebb49fe1d3b41aaf7bee77984ac4648e632cddc0"
	seedBytes, err := hex.DecodeString(seed)
	if err != nil {
		panic(err)
	}
	sender := "0xc0ee0c49b3be532975fdb2c02a3ae8dea70b58f53879f70ba256974627e23ee3"
	recipent := "0xa2909dd355aeaffd520847a35da1ef529c4cebd348bb4d69f00f32315118106e"
	objectIds, err := client.GetAllCoinObjectIds(types.SuiCoinType, sender)
	if err != nil {
		panic(err)
	}
	gasBudget := "4000000"
	amount := "4000000"
	suiObjectId := objectIds[0]

	tx, err := client.TransferSui(seedBytes, sender, recipent, suiObjectId, amount, gasBudget)
	if err != nil {
		panic(err)
	}
	fmt.Println(tx.Digest)
	//FgnntShXMos58nGvjsR1ZtjgXHunR1PmGbhBJ1HP9J7U
}

func TestSuiClient_TransferObject(t *testing.T) {
	seed := "dcee1121d8e62b67f9a5e8a9ebb49fe1d3b41aaf7bee77984ac4648e632cddc0"
	seedBytes, err := hex.DecodeString(seed)
	if err != nil {
		panic(err)
	}
	sender := "0xc0ee0c49b3be532975fdb2c02a3ae8dea70b58f53879f70ba256974627e23ee3"
	recipent := "0xa2909dd355aeaffd520847a35da1ef529c4cebd348bb4d69f00f32315118106e"
	objectIds, err := client.GetAllCoinObjectIds(types.SuiCoinType, sender)
	if err != nil {
		panic(err)
	}
	gasBudget := "4000000"
	suiObjectId := objectIds[0]
	gasObjectId := objectIds[1]

	signAndSubmitTx, err := client.TransferObject(seedBytes, sender, recipent, suiObjectId, gasObjectId, gasBudget)
	if err != nil {
		panic(err)
	}
	fmt.Println(signAndSubmitTx.Digest)
	//Cp7XhffS9uoenKRnavGuXUiF4zrhqMNSxrMmD8YonuG5
}

//func TestSuiClient_MoveCall(t *testing.T) {
//	seed := "491f69c043831b3b56ed8bf6e31080ad60bb3b44c5c41a6a978c37577634ba6b"
//	seedBytes, err := hex.DecodeString(seed)
//	if err != nil {
//		panic(err)
//	}
//	sender := "0x9e598d45b7ad757402e90d155c4e900045d30d21"
//	gasBudget := uint64(99)
//	_, objectList, err := client.Balance(types.SuiCoinType, sender)
//	if err != nil {
//		panic(err)
//	}
//	newObjectIds := BubbleSort(objectList)
//	gasObjectId := newObjectIds[0]
//
//	packageObjectId := ""
//	module := ""
//	funcation := ""
//	typeArguments := []string{}
//	arguemts := []string{}
//	signAndSubmitTx, err := client.MoveCall(seedBytes, sender, packageObjectId, module, funcation, gasObjectId, typeArguments, arguemts, gasBudget)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(signAndSubmitTx.EffectsCert)
//}

func TestSuiClient_Balance(t *testing.T) {
	//0xa2909dd355aeaffd520847a35da1ef529c4cebd348bb4d69f00f32315118106e
	balance, err := client.Balance(types.SuiCoinType, "0xa2909dd355aeaffd520847a35da1ef529c4cebd348bb4d69f00f32315118106e")
	if err != nil {
		panic(err)
	}
	fmt.Println(balance.String())
}

func TestDemo(t *testing.T) {
	data, err := base64.StdEncoding.DecodeString("AGLsaLe6fSvGG/YgrxirjhKqE21kVCcveOW9h0IiCZ1Ei/oAOmu95EnKjoBhLHcS2/2Ga2Ljw0BVnGrY6reYkwVDij1TvBYKLcfLNo8fq6GASb9yfo6uvuwNUBGkTf54wQ==")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}

func TestGenerateAddr(t *testing.T) {
	seed, err := crypto.NewRandSeed()
	if err != nil {
		panic(err)
	}
	fmt.Printf("seed:%x\n", seed)
	keyPair, err := crypto.NewKeyPairFromSeed(seed)
	if err != nil {
		panic(err)
	}
	fmt.Printf("privateKey:%x\n", keyPair.PrivateKey)
	fmt.Printf("publicKey:%x\n", keyPair.PublicKey)
	fmt.Printf("addess:%v\n", keyPair.Address())

}
