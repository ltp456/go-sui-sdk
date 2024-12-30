package go_sui_sdk

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ltp456/go-sui-sdk/crypto"
	"github.com/ltp456/go-sui-sdk/types"
	"golang.org/x/crypto/blake2b"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"reflect"
	"strconv"
)

type SuiClient struct {
	imp      *http.Client
	endpoint string
	debug    bool
}

func NewSuiClient(endpoint string) (*SuiClient, error) {
	client := &SuiClient{
		endpoint: endpoint,
		imp:      http.DefaultClient,
		debug:    false,
	}
	return client, nil
}

func (si *SuiClient) GetLatestCheckpointSequenceNumber() (uint64, error) {
	var result string
	err := si.post("sui_getLatestCheckpointSequenceNumber", nil, &result)
	if err != nil {
		return 0, err
	}
	latestPoint, err := strconv.ParseUint(result, 10, 64)
	if err != nil {
		return 0, err
	}
	return latestPoint, err
}

func (si *SuiClient) GetCheckpoint(point uint64) (*types.CheckPoint, error) {
	result := &types.CheckPoint{}
	params := Params{}
	params.AddValue(fmt.Sprintf("%v", point))
	err := si.post("sui_getCheckpoint", params, result)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (si *SuiClient) GetCheckpoints(start, limit uint64) (*types.CheckPoints, error) {
	result := &types.CheckPoints{}
	params := Params{}
	params.AddValue(fmt.Sprintf("%v", start))
	params.AddValue(limit)
	params.AddValue(false)
	err := si.post("sui_getCheckpoints", params, result)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (si *SuiClient) GetEvents(digest string) ([]types.TxEvent, error) {
	var result []types.TxEvent
	params := Params{}
	params.AddValue(digest)
	err := si.post("sui_getEvents", params, &result)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (si *SuiClient) GetTransactionBlock(digest string) (*types.TransactionBlock, error) {
	result := &types.TransactionBlock{}
	childParam := si.getDefaultTxOption()
	params := Params{}
	params.AddValue(digest)
	params.AddValue(childParam)
	err := si.post("sui_getTransactionBlock", params, result)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (si *SuiClient) TransactionsV1(checkpoint uint64) ([]types.Tx, error) {
	checkPoints, err := si.GetCheckpoint(checkpoint)
	if err != nil {
		return nil, err
	}
	var result []types.Tx
	for _, digest := range checkPoints.Transactions {
		transaction, err := si.GetTransactionBlock(digest)
		if err != nil {
			return nil, err
		}
		txList, err := transaction.Parse()
		if err != nil {
			continue
		}
		result = append(result, txList...)

	}
	return result, nil
}

func (si *SuiClient) Transactions(start, limit uint64) ([]types.Tx, uint64, error) {
	checkPoints, err := si.GetCheckpoints(start, limit)
	if err != nil {
		return nil, 0, err
	}

	var result []types.Tx
	for _, checkPoint := range checkPoints.Data {
		for _, tx := range checkPoint.Transactions {
			transaction, err := si.GetTransactionBlock(tx)
			if err != nil {
				return nil, 0, err
			}
			txList, err := transaction.Parse()
			if err != nil {
				continue
			}
			result = append(result, txList...)
		}
	}
	return result, uint64(len(checkPoints.Data)), err
}

func (si *SuiClient) GetReferenceGasPrice() (*big.Int, error) {
	var result string
	err := si.post("suix_getReferenceGasPrice", nil, &result)
	if err != nil {
		return nil, err
	}
	priceBig, ok := big.NewInt(0).SetString(result, 10)
	if !ok {
		return nil, fmt.Errorf("parse big error: %v", result)
	}
	return priceBig, err
}

func (si *SuiClient) GetAllBalances(address string) ([]types.AllBalance, error) {
	var result []types.AllBalance
	params := Params{}
	params.AddValue(address)
	err := si.post("suix_getAllBalances", params, &result)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (si *SuiClient) devInspectTransactionBlock(sender string, txBytes string) (*types.TransactionBlock, error) {
	result := &types.TransactionBlock{}
	params := Params{}
	params.AddValue(sender)
	params.AddValue(txBytes)
	params.AddValue("1000")
	err := si.post("sui_devInspectTransactionBlock", params, result)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (si *SuiClient) dryRunTransactionBlock(txBytes string) (*types.TransactionBlock, error) {
	result := &types.TransactionBlock{}
	params := Params{}
	params.AddValue(txBytes)
	err := si.post("sui_dryRunTransactionBlock", params, result)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (si *SuiClient) GetBalance(coinType types.CoinType, address string) (types.AllBalance, error) {
	result := types.AllBalance{}
	params := Params{}
	params.AddValue(address)
	params.AddValue(coinType)
	err := si.post("suix_getBalance", params, &result)
	if err != nil {
		return result, err
	}
	return result, err
}

func (si *SuiClient) Balance(coinType types.CoinType, address string) (*big.Int, error) {
	allBalances, err := si.GetAllBalances(address)
	if err != nil {
		return nil, err
	}
	for _, item := range allBalances {
		if types.CoinType(item.CoinType) == coinType {
			balanceBig, ok := big.NewInt(0).SetString(item.TotalBalance, 10)
			if !ok {
				return nil, fmt.Errorf("parse big error: %v", item.TotalBalance)
			}
			return balanceBig, nil
		}
	}
	return big.NewInt(0), nil
}

func (si *SuiClient) GetAllObjectId(address string) ([]string, error) {
	var objectIds []string
	var ownedObjects *types.ObjectInfo
	var err error
	var hasNetxPage bool
	ownedObjects, err = si.GetOwnedObjects("", address)
	if err != nil {
		return nil, err
	}
	for _, item := range ownedObjects.Data {
		objectIds = append(objectIds, item.Data.ObjectID)
	}

	hasNetxPage = ownedObjects.HasNextPage
	for hasNetxPage {
		ownedObjects, err = si.GetOwnedObjects(ownedObjects.NextCursor, address)
		if err != nil {
			return nil, err
		}
		for _, item := range ownedObjects.Data {
			objectIds = append(objectIds, item.Data.ObjectID)
		}
		hasNetxPage = ownedObjects.HasNextPage
	}
	return objectIds, nil

}

func (si *SuiClient) GetAllCoinObjectIds(coinType types.CoinType, address string) ([]string, error) {
	var objectIds []string
	var coinObjects *types.CoinObj
	var err error
	var hasNetxPage bool
	coinObjects, err = si.GetCoins(coinType, address, "")
	if err != nil {
		return nil, err
	}
	for _, item := range coinObjects.Data {
		if types.CoinType(item.CoinType) == coinType {
			objectIds = append(objectIds, item.CoinObjectID)
		}
	}

	hasNetxPage = coinObjects.HasNextPage
	for hasNetxPage {
		coinObjects, err = si.GetCoins(coinType, address, coinObjects.NextCursor)
		if err != nil {
			return nil, err
		}
		for _, item := range coinObjects.Data {
			if types.CoinType(item.CoinType) == coinType {
				objectIds = append(objectIds, item.CoinObjectID)
			}
		}
		hasNetxPage = coinObjects.HasNextPage
	}
	return objectIds, nil

}

func (si *SuiClient) GetCoins(coinType types.CoinType, address, cursor string) (*types.CoinObj, error) {
	result := &types.CoinObj{}
	params := Params{}
	params.AddValue(address)
	params.AddValue(coinType)
	if cursor != "" {
		params.AddValue(cursor)
	}
	err := si.post("suix_getCoins", params, result)
	return result, err
}

func (si *SuiClient) GetOwnedObjects(cursor, address string) (*types.ObjectInfo, error) {
	result := &types.ObjectInfo{}
	params := Params{}
	params.AddValue(address)
	if cursor != "" {
		params.AddValue(cursor)
	}
	err := si.post("suix_getOwnedObjects", params, result)
	return result, err
}

func (si *SuiClient) GetObject(objId string) (*types.ObjData, error) {
	result := &types.ObjData{}
	params := Params{}
	params.AddValue(objId)
	err := si.post("sui_getObject", params, result)
	return result, err
}

func (si *SuiClient) DryRunTransactionBlock(coinType types.CoinType, payAllSui bool, sender string, objectIds []string, recipient string, amount string, gasBudget string) (*types.TransactionBlock, error) {
	var unsignedTx *types.UnsignedTx
	var err error
	if coinType == types.SuiCoinType {
		if payAllSui {
			unsignedTx, err = si.payAllSui(sender, recipient, objectIds, gasBudget)
			if err != nil {
				return nil, err
			}
		} else {
			unsignedTx, err = si.paySui(sender, objectIds, []string{recipient}, []string{amount}, gasBudget)
			if err != nil {
				return nil, err
			}
		}

	} else {
		unsignedTx, err = si.pay(sender, "", nil, []string{recipient}, []string{amount}, gasBudget)
		if err != nil {
			return nil, err
		}
	}
	transactionBlock, err := si.dryRunTransactionBlock(unsignedTx.TxBytes)
	if err != nil {
		return nil, err
	}
	return transactionBlock, nil

}

func (si *SuiClient) DevInspectTransactionBlock(coinType types.CoinType, payAllSui bool, sender string, objectIds []string, recipient string, amount string, gasBudget string) (*types.TransactionBlock, error) {
	var unsignedTx *types.UnsignedTx
	var err error
	if coinType == types.SuiCoinType {
		if payAllSui {
			unsignedTx, err = si.payAllSui(sender, recipient, objectIds, gasBudget)
			if err != nil {
				return nil, err
			}
		} else {
			unsignedTx, err = si.paySui(sender, objectIds, []string{recipient}, []string{amount}, gasBudget)
			if err != nil {
				return nil, err
			}
		}

	} else {
		unsignedTx, err = si.pay(sender, "", nil, []string{recipient}, []string{amount}, gasBudget)
		if err != nil {
			return nil, err
		}
	}
	transactionBlock, err := si.devInspectTransactionBlock(sender, unsignedTx.TxBytes)
	if err != nil {
		return nil, err
	}
	return transactionBlock, nil

}

func (si *SuiClient) Transfer(coinType types.CoinType, seed []byte, sender string, allObjectIds, recipient []string, amount []string, gasBudget string) (*types.TransactionBlock, error) {
	var unsignedTx *types.UnsignedTx
	var err error
	if coinType == types.SuiCoinType {
		unsignedTx, err = si.paySui(sender, allObjectIds, recipient, amount, gasBudget)
		if err != nil {
			return nil, err
		}
	} else {
		unsignedTx, err = si.pay(sender, "", nil, recipient, amount, gasBudget)
		if err != nil {
			return nil, err
		}
	}
	result, err := si.signAndSubmitTx(seed, unsignedTx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (si *SuiClient) Pay(seed []byte, sender, gasObjectId string, inputCoins, recipient []string, amount []string, gasBudget string) (*types.TransactionBlock, error) {
	unsignedTx, err := si.pay(sender, gasObjectId, inputCoins, recipient, amount, gasBudget)
	if err != nil {
		return nil, err
	}
	submitTx, err := si.signAndSubmitTx(seed, unsignedTx)
	if err != nil {
		return nil, err
	}
	return submitTx, nil
}

func (si *SuiClient) PayAllSui(seed []byte, sender, recipient string, suiObjectId []string, gasBudget string) (*types.TransactionBlock, error) {
	unsignedTx, err := si.payAllSui(sender, recipient, suiObjectId, gasBudget)
	if err != nil {
		return nil, err
	}
	submitTx, err := si.signAndSubmitTx(seed, unsignedTx)
	if err != nil {
		return nil, err
	}
	return submitTx, nil
}

func (si *SuiClient) PaySui(seed []byte, sender string, inputCoins, recipient []string, amount []string, gasBudget string) (*types.TransactionBlock, error) {
	unsignedTx, err := si.paySui(sender, inputCoins, recipient, amount, gasBudget)
	if err != nil {
		return nil, err
	}
	submitTx, err := si.signAndSubmitTx(seed, unsignedTx)
	if err != nil {
		return nil, err
	}
	return submitTx, nil
}

func (si *SuiClient) TransferSui(seed []byte, sender, recipient, suiObjectId string, amount string, gasBudget string) (*types.TransactionBlock, error) {
	unsignedTx, err := si.transferSui(sender, recipient, suiObjectId, amount, gasBudget)
	if err != nil {
		return nil, err
	}
	submitTx, err := si.signAndSubmitTx(seed, unsignedTx)
	if err != nil {
		return nil, err
	}
	return submitTx, nil
}

func (si *SuiClient) TransferObject(seed []byte, sender, recipient, suiObjectId, gasObjectId string, gasBudget string) (*types.TransactionBlock, error) {
	unsignedTx, err := si.transferObject(sender, recipient, suiObjectId, gasObjectId, gasBudget)
	if err != nil {
		return nil, err
	}
	submitTx, err := si.signAndSubmitTx(seed, unsignedTx)
	if err != nil {
		return nil, err
	}
	return submitTx, nil
}

func (si *SuiClient) MoveCall(seed []byte, sender, packageObjectId, module, function, gasObjectId string, typeArguments, arguments []string, gasBudget uint64) (*types.TransactionBlock, error) {
	unsignedTx, err := si.moveCall(sender, packageObjectId, module, function, gasObjectId, typeArguments, arguments, gasBudget)
	if err != nil {
		return nil, err
	}
	submitTx, err := si.signAndSubmitTx(seed, unsignedTx)
	if err != nil {
		return nil, err
	}
	return submitTx, nil
}

func (si *SuiClient) signAndSubmitTx(seed []byte, unsignedTx *types.UnsignedTx) (*types.TransactionBlock, error) {
	keyPair, err := crypto.NewKeyPairFromSeed(seed)
	if err != nil {
		return nil, err
	}
	txBytes, err := base64.StdEncoding.DecodeString(unsignedTx.TxBytes)
	if err != nil {
		return nil, err
	}

	txData := make([]byte, 0)
	txData = append(txData, types.IntentFlag...)
	txData = append(txData, txBytes...)
	txHash := blake2b.Sum256(txData)
	signature, err := keyPair.Sign(txHash[:])
	if err != nil {
		return nil, err
	}

	signatureData := make([]byte, 0)
	signatureData = append(signatureData, []byte{0}...)
	signatureData = append(signatureData, signature...)
	signatureData = append(signatureData, keyPair.PublicKey...)

	base64Signature := base64.StdEncoding.EncodeToString(signatureData)
	if err != nil {
		return nil, err
	}
	transaction, err := si.ExecuteTransactionBlock(unsignedTx.TxBytes, base64Signature, types.WaitForLocalExecution.String(), si.getDefaultTxOption())
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (si *SuiClient) getDefaultTxOption() MapParams {
	params := MapParams{}
	params.SetKey("showInput", true)
	params.SetKey("showRawInput", false)
	params.SetKey("showEffects", true)
	params.SetKey("showEvents", false)
	params.SetKey("showObjectChanges", false)
	params.SetKey("showBalanceChanges", true)
	return params
}

func (si *SuiClient) ExecuteTransactionBlock(txBytes, signature, requestType string, options MapParams) (*types.TransactionBlock, error) {
	result := &types.TransactionBlock{}
	params := Params{}
	params.AddValue(txBytes)
	params.AddValue([]string{signature})
	params.AddValue(options)
	params.AddValue(requestType)
	err := si.post("sui_executeTransactionBlock", params, result)
	return result, err
}

func (si *SuiClient) pay(sender, gasObjectId string, inputCoins, recipient []string, amount []string, gasBudget string) (*types.UnsignedTx, error) {
	result := &types.UnsignedTx{}
	params := Params{}
	params.AddValue(sender)
	params.AddValue(inputCoins)
	params.AddValue(recipient)
	params.AddValue(amount)
	params.AddValue(gasObjectId)
	params.AddValue(gasBudget)
	err := si.post("unsafe_pay", params, result)
	return result, err
}

func (si *SuiClient) paySui(sender string, inputCoins, recipient []string, amount []string, gasBudget string) (*types.UnsignedTx, error) {
	result := &types.UnsignedTx{}
	params := Params{}
	params.AddValue(sender)
	params.AddValue(inputCoins)
	params.AddValue(recipient)
	params.AddValue(amount)
	params.AddValue(gasBudget)
	err := si.post("unsafe_paySui", params, result)
	return result, err
}

func (si *SuiClient) payAllSui(sender, recipient string, inputCoins []string, gasBudget string) (*types.UnsignedTx, error) {
	result := &types.UnsignedTx{}
	params := Params{}
	params.AddValue(sender)
	params.AddValue(inputCoins)
	params.AddValue(recipient)
	params.AddValue(gasBudget)
	err := si.post("unsafe_payAllSui", params, result)
	return result, err
}

func (si *SuiClient) transferObject(sender, recipient, objectId, gasObjectId string, gasBudget string) (*types.UnsignedTx, error) {
	result := &types.UnsignedTx{}
	params := Params{}
	params.AddValue(sender)
	params.AddValue(objectId)
	params.AddValue(gasObjectId)
	params.AddValue(gasBudget)
	params.AddValue(recipient)
	err := si.post("unsafe_transferObject", params, result)
	return result, err
}

func (si *SuiClient) transferSui(sender, recipient, suiObjectId string, amount, gasBudget string) (*types.UnsignedTx, error) {
	result := &types.UnsignedTx{}
	params := Params{}
	params.AddValue(sender)
	params.AddValue(suiObjectId)
	params.AddValue(gasBudget)
	params.AddValue(recipient)
	params.AddValue(amount)
	err := si.post("unsafe_transferSui", params, result)
	return result, err
}

func (si *SuiClient) moveCall(sender, packageObjectId, module, function, gasObjectId string, typeArguments, arguments []string, gasBudget uint64) (*types.UnsignedTx, error) {
	result := &types.UnsignedTx{}
	params := Params{}
	params.AddValue(sender)
	params.AddValue(packageObjectId)
	params.AddValue(module)
	params.AddValue(function)
	params.AddValue(gasObjectId)
	params.AddValue(typeArguments)
	params.AddValue(arguments)
	params.AddValue(gasBudget)
	err := si.post("sui_moveCall", params, result)
	return result, err
}

func (si *SuiClient) GetTotalTransactionBlocks() (uint64, error) {
	var result string
	err := si.post("sui_getTotalTransactionBlocks", nil, &result)
	number, err := strconv.ParseUint(result, 10, 64)
	if err != nil {
		return 0, err
	}
	return number, err
}

func (si *SuiClient) post(method string, param Params, value interface{}, options ...Option) error {
	return si.httpReq(http.MethodPost, method, param, value, options...)
}

func (si *SuiClient) newRequest(httpMethod, url, method string, param interface{}) (*http.Request, int, error) {
	jsonRpc := NewJsonRpc(method, param)
	reqData, err := json.Marshal(jsonRpc)
	if err != nil {
		return nil, 0, err
	}
	req, err := http.NewRequest(httpMethod, url, bytes.NewReader(reqData))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, jsonRpc.ID, nil
}

func (si *SuiClient) httpReq(httpMethod, method string, param Params, value interface{}, options ...Option) (err error) {
	vi := reflect.ValueOf(value)
	if vi.Kind() != reflect.Ptr {
		return fmt.Errorf("value must be pointer")
	}

	req, id, err := si.newRequest(httpMethod, si.endpoint, method, param)
	if err != nil {
		return err
	}
	if si.debug {
		if param != nil {
			requestData, err := json.Marshal(param)
			if err != nil {
				return fmt.Errorf("%v", err)
			}
			log.Printf("httpReq request: %v  %v \n", method, string(requestData))
		}
	}
	resp, err := si.imp.Do(req)
	if err != nil {
		if resp != nil && resp.Body != nil {
			defer resp.Body.Close()
		}

		return err
	}
	if resp == nil || resp.StatusCode < http.StatusOK || resp.StatusCode > 300 {
		data, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("response err: %v %v %v", resp.StatusCode, resp.Status, string(data))
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if si.debug {
		log.Printf("httpReq response: %v %v \n", method, string(data))
	}
	jsonResp := JsonResp{}
	err = json.Unmarshal(data, &jsonResp)
	if err != nil {
		return err
	}
	if id != jsonResp.Id {
		return fmt.Errorf("%v jsonRpc reqId  %v not match RespId %v", method, id, jsonResp.Id)
	}
	if jsonResp.Error.Code != 0 {
		return fmt.Errorf("jsonRpc error: %v", jsonResp.Error.Message)
	}

	err = Unmarshal(jsonResp.Result, value)
	if err != nil {
		return fmt.Errorf("unmarashal error: %v %v", err, string(jsonResp.Result[:100]))
	}
	return nil

}
