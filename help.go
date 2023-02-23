package go_sui_sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-sui-sdk/types"
	"math/big"
	"reflect"
	"strings"
)

func GetMatchObjectIds(coinType types.CoinResType, amount uint64, objList []*types.Object) ([]string, error) {
	var objectIds []string
	amountBig := big.NewInt(0).SetUint64(amount)
	totalBig := big.NewInt(0)
	for _, object := range objList {
		if object.Details.Data.Type != coinType.String() {
			continue
		}
		if totalBig.Cmp(amountBig) >= 0 {
			break
		}
		totalBig = big.NewInt(0).Add(totalBig, object.Details.Data.Fields.Balance)
		objectIds = append(objectIds, object.ObjectId)
	}
	return objectIds, nil
}

func BubbleSort(res []*types.Object) []string {
	var objectIds []string
	flag := true
	for i := 0; i < len(res) && flag; i++ {
		flag = false
		for j := len(res) - 1; j > i; j-- {
			if res[j-1].Details.Data.Fields.Balance.Cmp(res[j].Details.Data.Fields.Balance) < 0 {
				res[j-1], res[j] = res[j], res[j-1]
				flag = true
			}
		}
	}
	for _, item := range res {
		objectIds = append(objectIds, item.ObjectId)
	}
	return objectIds
}

func Marshal(obj interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encode := json.NewEncoder(&buf)
	err := encode.Encode(obj)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func Unmarshal(data []byte, value interface{}) error {
	typeOf := reflect.ValueOf(value)
	if typeOf.Kind() != reflect.Ptr {
		return fmt.Errorf("value is mutst pointer")
	}
	decode := json.NewDecoder(strings.NewReader(string(data)))
	decode.UseNumber()
	err := decode.Decode(value)
	if err != nil {
		return err
	}
	return err
}
