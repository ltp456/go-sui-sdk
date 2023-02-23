package go_sui_sdk

import (
	"encoding/hex"
	"fmt"
	"go-sui-sdk/crypto"
	"go-sui-sdk/types"
	"testing"
)

var client *SuiClient
var err error

func init() {
	endpoint := "https://fullnode.devnet.sui.io:443"
	//endpoint := "http://127.0.0.1:9000"
	client, err = NewSuiClient(endpoint)
	if err != nil {
		panic(err)
	}
}

func TestSuiClient_GetTotalTransactionNumber(t *testing.T) {
	number, err := client.GetTotalTransactionNumber()
	if err != nil {
		panic(err)
	}
	fmt.Println(number)
}

func TestSuiClient_GetTransaction(t *testing.T) {
	tx, err := client.GetTransaction("+pdagEyZE7XqJCsQ3uUYZH053HWJ69dON0hEEB96HUI=")
	if err != nil {
		panic(err)
	}
	fmt.Println(tx)
}

func TestSuiClient_GetTransactionsInRange(t *testing.T) {
	result, err := client.GetTransactionsInRange(3833520, 3833522)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

func TestSuiClient_Transactions(t *testing.T) {

	pageNum := 20
	start := 910457
	for i := start; i+pageNum < 1000000000000000; i = i + pageNum {
		txes, err := client.Transactions(uint64(i), uint64(i+pageNum))
		if err != nil {
			panic(err)
		}
		for _, tx := range txes {
			fmt.Println(tx.Height, tx.Hash, tx.Sender, tx.Recipient, tx.Amount.String(), tx.CoinType.String(), tx.Gas.String(), tx.TxMethod)
		}

	}
}

func TestSuiClient_ExecuteTransaction(t *testing.T) {
	transaction, err := client.ExecuteTransaction(
		"", "", "", "", "")
	if err != nil {
		panic(err)
	}
	fmt.Println(transaction)
}

func TestSuiClient_Transfer(t *testing.T) {
	seed := "491f69c043831b3b56ed8bf6e31080ad60bb3b44c5c41a6a978c37577634ba6b"
	seedBytes, err := hex.DecodeString(seed)
	if err != nil {
		panic(err)
	}
	sender := "0x9e598d45b7ad757402e90d155c4e900045d30d21"
	recipent := "0x7460b88ade060545505c548eb8f4269ec6d52492"
	//80158286
	//375239039
	amount := uint64(40)
	gasBudget := uint64(37)

	_, objectList, err := client.Balance(types.SuiCoinType, sender)
	if err != nil {
		panic(err)
	}

	//objectIds, err := GetMatchObjectIds(types.SuiCoinType, amount, objectList)
	//if err != nil {
	//	panic(err)
	//}
	newObjectIds := BubbleSort(objectList)
	gasObjectId := "0x7d7df1aa4bde636f87790a1f58f45a9af4aacfd8"

	tx, err := client.Transfer(types.SuiCoinType, seedBytes, sender, gasObjectId, []string{recipent}, newObjectIds, []uint64{amount}, gasBudget)
	if err != nil {
		panic(err)
	}
	fmt.Println(tx.EffectsCert.Certificate.TransactionDigest)
}

func TestSuiClient_Pay(t *testing.T) {

	seed := "491f69c043831b3b56ed8bf6e31080ad60bb3b44c5c41a6a978c37577634ba6b"
	seedBytes, err := hex.DecodeString(seed)
	if err != nil {
		panic(err)
	}
	sender := "0x9e598d45b7ad757402e90d155c4e900045d30d21"
	recipent := "0x7460b88ade060545505c548eb8f4269ec6d52492"
	amount := uint64(5555)
	gasBudget := uint64(99)

	_, objectList, err := client.Balance(types.SuiCoinType, sender)
	if err != nil {
		panic(err)
	}
	newObjectIds := BubbleSort(objectList)
	gasObjectId := newObjectIds[0]
	inputObjectIds := newObjectIds[1:]
	tx, err := client.Pay(seedBytes, sender, gasObjectId, inputObjectIds, []string{recipent}, []uint64{amount}, gasBudget)
	if err != nil {
		panic(err)
	}
	fmt.Println(tx.EffectsCert.Certificate.TransactionDigest)

}

func TestSuiClient_PaySui(t *testing.T) {

	seed := "491f69c043831b3b56ed8bf6e31080ad60bb3b44c5c41a6a978c37577634ba6b"
	seedBytes, err := hex.DecodeString(seed)
	if err != nil {
		panic(err)
	}
	sender := "0x9e598d45b7ad757402e90d155c4e900045d30d21"
	recipent := "0x7460b88ade060545505c548eb8f4269ec6d52492"
	gasBudget := uint64(200)
	amount := uint64(269994304)
	_, objectList, err := client.Balance(types.SuiCoinType, sender)
	if err != nil {
		panic(err)
	}

	newObjectIds := BubbleSort(objectList)
	tx, err := client.PaySui(seedBytes, sender, newObjectIds, []string{recipent}, []uint64{amount}, gasBudget)
	if err != nil {
		panic(err)
	}
	fmt.Println(tx.EffectsCert.Certificate.TransactionDigest)
}

func TestSuiClient_PayAllSui(t *testing.T) {
	seed := "491f69c043831b3b56ed8bf6e31080ad60bb3b44c5c41a6a978c37577634ba6b"
	seedBytes, err := hex.DecodeString(seed)
	if err != nil {
		panic(err)
	}
	sender := "0x9e598d45b7ad757402e90d155c4e900045d30d21"
	recipent := "0x7460b88ade060545505c548eb8f4269ec6d52492"
	gasBudget := uint64(200)
	_, objectList, err := client.Balance(types.SuiCoinType, sender)
	if err != nil {
		panic(err)
	}

	newObjectIds := BubbleSort(objectList)

	tx, err := client.PayAllSui(seedBytes, sender, recipent, newObjectIds, gasBudget)
	if err != nil {
		panic(err)
	}
	fmt.Println(tx.EffectsCert.Certificate.TransactionDigest)
}

func TestSuiClient_TransferSui(t *testing.T) {

	seed := "491f69c043831b3b56ed8bf6e31080ad60bb3b44c5c41a6a978c37577634ba6b"
	seedBytes, err := hex.DecodeString(seed)
	if err != nil {
		panic(err)
	}
	sender := "0x9e598d45b7ad757402e90d155c4e900045d30d21"
	recipent := "0x7460b88ade060545505c548eb8f4269ec6d52492"
	gasBudget := uint64(99)
	amount := uint64(123456)
	_, objectList, err := client.Balance(types.SuiCoinType, sender)
	if err != nil {
		panic(err)
	}
	inputCoinObjectIds := BubbleSort(objectList)
	suiObjectId := inputCoinObjectIds[0]
	tx, err := client.TransferSui(seedBytes, sender, recipent, suiObjectId, amount, gasBudget)
	if err != nil {
		panic(err)
	}
	fmt.Println(tx.EffectsCert.Certificate.TransactionDigest)
}

func TestSuiClient_TransferObject(t *testing.T) {
	seed := "491f69c043831b3b56ed8bf6e31080ad60bb3b44c5c41a6a978c37577634ba6b"
	seedBytes, err := hex.DecodeString(seed)
	if err != nil {
		panic(err)
	}
	sender := "0x9e598d45b7ad757402e90d155c4e900045d30d21"
	recipent := "0x7460b88ade060545505c548eb8f4269ec6d52492"
	gasBudget := uint64(99)
	_, objectList, err := client.Balance(types.SuiCoinType, sender)
	if err != nil {
		panic(err)
	}
	newObjectIds := BubbleSort(objectList)
	gasObjectId := newObjectIds[0]
	inputObjectIds := newObjectIds[1]
	submitTx, err := client.TransferObject(seedBytes, sender, recipent, inputObjectIds, gasObjectId, gasBudget)
	if err != nil {
		panic(err)
	}
	fmt.Println(submitTx.EffectsCert.Certificate.TransactionDigest)
}

func TestSuiClient_MoveCall(t *testing.T) {
	seed := "491f69c043831b3b56ed8bf6e31080ad60bb3b44c5c41a6a978c37577634ba6b"
	seedBytes, err := hex.DecodeString(seed)
	if err != nil {
		panic(err)
	}
	sender := "0x9e598d45b7ad757402e90d155c4e900045d30d21"
	gasBudget := uint64(99)
	_, objectList, err := client.Balance(types.SuiCoinType, sender)
	if err != nil {
		panic(err)
	}
	newObjectIds := BubbleSort(objectList)
	gasObjectId := newObjectIds[0]

	packageObjectId := ""
	module := ""
	funcation := ""
	typeArguments := []string{}
	arguemts := []string{}
	submitTx, err := client.MoveCall(seedBytes, sender, packageObjectId, module, funcation, gasObjectId, typeArguments, arguemts, gasBudget)
	if err != nil {
		panic(err)
	}
	fmt.Println(submitTx.EffectsCert.Certificate.TransactionDigest)
}

func TestSuiClient_GetObject(t *testing.T) {
	object, err := client.GetObject("0x7d7df1aa4bde636f87790a1f58f45a9af4aacfd8")
	if err != nil {
		panic(err)
	}
	fmt.Println(object)
}

func TestSuiClient_GetObjectsOwnedByAddress(t *testing.T) {
	result, err := client.GetObjectsOwnedByAddress("0x7460b88ade060545505c548eb8f4269ec6d52492")
	if err != nil {
		panic(err)
	}
	for _, item := range result {
		fmt.Println(item.ObjectID, item)
	}
}

func TestSuiClient_Balance(t *testing.T) {
	balance, _, err := client.Balance(types.SuiCoinType, "0x9e598d45b7ad757402e90d155c4e900045d30d21")
	if err != nil {
		panic(err)
	}
	fmt.Println(balance.String())
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
