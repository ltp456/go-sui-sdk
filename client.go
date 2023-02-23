package go_sui_sdk

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-sui-sdk/crypto"
	"go-sui-sdk/types"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"reflect"
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

func (si *SuiClient) Balance(coinType types.CoinType, address string) (*big.Int, []*types.Object, error) {
	objectInfos, err := si.GetObjectsOwnedByAddress(address)
	if err != nil {
		return nil, nil, err
	}
	total := big.NewInt(0)
	var objectIds []*types.Object
	for _, object := range objectInfos {
		if object.Type == coinType.String() {
			objDetail, err := si.GetObject(object.ObjectID)
			if err != nil {
				return nil, nil, err
			}
			if objDetail.Status != types.Exists {
				continue
			}
			if objDetail.Details.Data.Type != coinType.String() {
				continue
			}
			total = big.NewInt(0).Add(total, objDetail.Details.Data.Fields.Balance)
			objDetail.ObjectId = object.ObjectID
			objectIds = append(objectIds, objDetail)

		}
	}
	return total, objectIds, nil
}

func (si *SuiClient) GetObjectsOwnedByAddress(address string) ([]types.ObjectInfo, error) {
	var result []types.ObjectInfo
	params := Params{}
	params.AddValue(address)
	err := si.post("sui_getObjectsOwnedByAddress", params, &result)
	return result, err
}

func (si *SuiClient) GetObject(objId string) (*types.Object, error) {
	result := &types.Object{}
	params := Params{}
	params.AddValue(objId)
	err := si.post("sui_getObject", params, &result)
	return result, err
}

func (si *SuiClient) Transfer(coinType types.CoinType, seed []byte, sender, gasObjectId string, recipient, objectIds []string, amount []uint64, gasBudget uint64) (*types.SubmitTx, error) {
	if len(recipient) != len(amount) {
		return nil, fmt.Errorf("recipient lenght  %v not equal amount length %v", len(recipient), len(amount))
	}
	var unsignedTx *types.UnsignedTx
	var err error
	if coinType == types.SuiCoinType {
		unsignedTx, err = si.paySui(sender, objectIds, recipient, amount, gasBudget)
		if err != nil {
			return nil, err
		}
	} else {
		unsignedTx, err = si.pay(sender, gasObjectId, objectIds, recipient, amount, gasBudget)
		if err != nil {
			return nil, err
		}
	}
	result, err := si.submitTx(seed, unsignedTx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (si *SuiClient) Pay(seed []byte, sender, gasObjectId string, inputCoins, recipient []string, amount []uint64, gasBudget uint64) (*types.SubmitTx, error) {
	unsignedTx, err := si.pay(sender, gasObjectId, inputCoins, recipient, amount, gasBudget)
	if err != nil {
		return nil, err
	}
	submitTx, err := si.submitTx(seed, unsignedTx)
	if err != nil {
		return nil, err
	}
	return submitTx, nil
}

func (si *SuiClient) PayAllSui(seed []byte, sender, recipient string, suiObjectId []string, gasBudget uint64) (*types.SubmitTx, error) {
	unsignedTx, err := si.payAllSui(sender, recipient, suiObjectId, gasBudget)
	if err != nil {
		return nil, err
	}
	submitTx, err := si.submitTx(seed, unsignedTx)
	if err != nil {
		return nil, err
	}
	return submitTx, nil
}

func (si *SuiClient) PaySui(seed []byte, sender string, inputCoins, recipient []string, amount []uint64, gasBudget uint64) (*types.SubmitTx, error) {
	unsignedTx, err := si.paySui(sender, inputCoins, recipient, amount, gasBudget)
	if err != nil {
		return nil, err
	}
	submitTx, err := si.submitTx(seed, unsignedTx)
	if err != nil {
		return nil, err
	}
	return submitTx, nil
}

func (si *SuiClient) TransferSui(seed []byte, sender, recipient, suiObjectId string, amount uint64, gasBudget uint64) (*types.SubmitTx, error) {
	unsignedTx, err := si.transferSui(sender, recipient, suiObjectId, amount, gasBudget)
	if err != nil {
		return nil, err
	}
	submitTx, err := si.submitTx(seed, unsignedTx)
	if err != nil {
		return nil, err
	}
	return submitTx, nil
}

func (si *SuiClient) TransferObject(seed []byte, sender, recipient, suiObjectId, gasObjectId string, gasBudget uint64) (*types.SubmitTx, error) {
	unsignedTx, err := si.transferObject(sender, recipient, suiObjectId, gasObjectId, gasBudget)
	if err != nil {
		return nil, err
	}
	submitTx, err := si.submitTx(seed, unsignedTx)
	if err != nil {
		return nil, err
	}
	return submitTx, nil
}

func (si *SuiClient) MoveCall(seed []byte, sender, packageObjectId, module, function, gasObjectId string, typeArguments, arguments []string, gasBudget uint64) (*types.SubmitTx, error) {
	unsignedTx, err := si.moveCall(sender, packageObjectId, module, function, gasObjectId, typeArguments, arguments, gasBudget)
	if err != nil {
		return nil, err
	}
	submitTx, err := si.submitTx(seed, unsignedTx)
	if err != nil {
		return nil, err
	}
	return submitTx, nil
}

func (si *SuiClient) submitTx(seed []byte, unsignedTx *types.UnsignedTx) (*types.SubmitTx, error) {
	keyPair, err := crypto.NewKeyPairFromSeed(seed)
	if err != nil {
		return nil, err
	}
	txBytes, err := base64.StdEncoding.DecodeString(unsignedTx.TxBytes)
	if err != nil {
		return nil, err
	}
	signature, err := keyPair.Sign(txBytes)
	if err != nil {
		return nil, err
	}
	base64Signature := base64.StdEncoding.EncodeToString(signature)
	publicKeyBase64 := base64.StdEncoding.EncodeToString(keyPair.PublicKey)

	transaction, err := si.ExecuteTransaction(unsignedTx.TxBytes, keyPair.Type().String(), base64Signature, publicKeyBase64, types.WaitForLocalExecution.String())
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (si *SuiClient) ExecuteTransaction(txBytes, sigScheme, signature, pubKey, requestType string) (*types.SubmitTx, error) {
	result := &types.SubmitTx{}
	params := Params{}
	params.AddValue(txBytes)
	params.AddValue(sigScheme)
	params.AddValue(signature)
	params.AddValue(pubKey)
	params.AddValue(requestType)
	err := si.post("sui_executeTransaction", params, result)
	return result, err
}

func (si *SuiClient) pay(sender, gasObjectId string, inputCoins, recipient []string, amount []uint64, gasBudget uint64) (*types.UnsignedTx, error) {
	result := &types.UnsignedTx{}
	params := Params{}
	params.AddValue(sender)
	params.AddValue(inputCoins)
	params.AddValue(recipient)
	params.AddValue(amount)
	params.AddValue(gasObjectId)
	params.AddValue(gasBudget)
	err := si.post("sui_pay", params, result)
	return result, err
}

func (si *SuiClient) paySui(sender string, inputCoins, recipient []string, amount []uint64, gasBudget uint64) (*types.UnsignedTx, error) {
	result := &types.UnsignedTx{}
	params := Params{}
	params.AddValue(sender)
	params.AddValue(inputCoins)
	params.AddValue(recipient)
	params.AddValue(amount)
	params.AddValue(gasBudget)
	err := si.post("sui_paySui", params, result)
	return result, err
}

func (si *SuiClient) payAllSui(sender, recipient string, inputCoins []string, gasBudget uint64) (*types.UnsignedTx, error) {
	result := &types.UnsignedTx{}
	params := Params{}
	params.AddValue(sender)
	params.AddValue(inputCoins)
	params.AddValue(recipient)
	params.AddValue(gasBudget)
	err := si.post("sui_payAllSui", params, result)
	return result, err
}

func (si *SuiClient) transferObject(sender, recipient, objectId, gasObjectId string, gasBudget uint64) (*types.UnsignedTx, error) {
	result := &types.UnsignedTx{}
	params := Params{}
	params.AddValue(sender)
	params.AddValue(objectId)
	params.AddValue(gasObjectId)
	params.AddValue(gasBudget)
	params.AddValue(recipient)
	err := si.post("sui_transferObject", params, result)
	return result, err
}

func (si *SuiClient) transferSui(sender, recipient, suiObjectId string, amount, gasBudget uint64) (*types.UnsignedTx, error) {
	result := &types.UnsignedTx{}
	params := Params{}
	params.AddValue(sender)
	params.AddValue(suiObjectId)
	params.AddValue(gasBudget)
	params.AddValue(recipient)
	params.AddValue(amount)
	err := si.post("sui_transferSui", params, result)
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

func (si *SuiClient) Transactions(start, end uint64) ([]types.Tx, error) {
	txDigestList, err := si.GetTransactionsInRange(start, end)
	if err != nil {
		return nil, err
	}
	var result []types.Tx
	for index, hash := range txDigestList {
		transaction, err := si.GetTransaction(hash)
		if err != nil {
			return nil, err
		}
		height := start + uint64(index)
		txes, err := transaction.Parse(height)
		if err != nil {
			continue
		}
		result = append(result, txes...)
	}
	return result, nil
}

func (si *SuiClient) GetTransactionsInRange(start, end uint64) ([]string, error) {
	var result []string
	params := Params{}
	params.AddValue(start)
	params.AddValue(end)
	err := si.post("sui_getTransactionsInRange", params, &result)
	return result, err
}

func (si *SuiClient) GetTransaction(digest string) (*types.Transaction, error) {
	result := &types.Transaction{}
	params := Params{}
	params.AddValue(digest)
	err := si.post("sui_getTransaction", params, result)
	return result, err
}

func (si *SuiClient) GetLedgerNumber() (uint64, error) {
	return si.GetTotalTransactionNumber()
}

func (si *SuiClient) GetTotalTransactionNumber() (uint64, error) {
	result := uint64(0)
	err := si.post("sui_getTotalTransactionNumber", nil, &result)
	return result, err
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
		return fmt.Errorf("unmarashal error: %v %v", err, string(jsonResp.Result))
	}
	return nil

}
