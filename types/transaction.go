package types

import (
	"fmt"
	"math/big"
	"strconv"
	"time"
)

type UnsignedTx struct {
	TxBytes      string      `json:"txBytes"`
	InputObjects interface{} `json:"inputObjects"`
	Gas          interface{} `json:"gas"`
}

type SubmitTx struct {
	EffectsCert Transaction `json:"EffectsCert"`
}

// --------------
type Tx struct {
	Sender         string
	Recipient      string
	Amount         *big.Int
	Gas            *big.Int
	Hash           string
	Checkpoint     uint64
	Epoch          string
	TxSignatures   []string
	MessageVersion string
	CoinType       CoinType
	Time           time.Time
}

type TransactionBlock struct {
	Digest      string      `json:"digest"`
	Transaction Transaction `json:"transaction"`
	//RawTransaction string           `json:"rawTransaction"`
	Effects Effects `json:"effects"`
	//Events         []interface{}    `json:"events"`
	//ObjectChanges  []ObjectChanges  `json:"objectChanges"`
	BalanceChanges []BalanceChanges `json:"balanceChanges"`
	TimestampMs    string           `json:"timestampMs"`
	Checkpoint     string           `json:"checkpoint"`
}

func (tb *TransactionBlock) Parse() ([]Tx, error) {
	if tb.Effects.Status.Status != TxSuccess {
		return nil, fmt.Errorf("tx fail: %v", tb.Effects.Status)
	}

	checkPoint, err := strconv.ParseUint(tb.Checkpoint, 10, 64)
	if err != nil {
		return nil, err
	}
	gasUsed, err := tb.GetGasUsed()
	if err != nil {
		return nil, err
	}
	timestamp, err := strconv.ParseInt(tb.TimestampMs, 10, 64)
	if err != nil {
		return nil, err
	}
	txTimestamp := time.UnixMilli(timestamp)
	var tmpList []BalanceChanges
	for _, balChange := range tb.BalanceChanges {
		amountBig, ok := big.NewInt(0).SetString(balChange.Amount, 10)
		if !ok {
			continue
		}
		if amountBig.Cmp(big.NewInt(0)) > 0 {
			tmpList = append(tmpList, balChange)
		}
	}
	var txList []Tx

	for _, balChange := range tmpList {
		amountBig, ok := big.NewInt(0).SetString(balChange.Amount, 10)
		if !ok {
			continue
		}
		tx := Tx{
			Hash:           tb.Digest,
			Checkpoint:     checkPoint,
			TxSignatures:   tb.Transaction.TxSignatures,
			Sender:         tb.Transaction.Data.Sender,
			Gas:            gasUsed,
			Epoch:          tb.Effects.ExecutedEpoch,
			MessageVersion: tb.Effects.MessageVersion,
			Recipient:      balChange.Owner.AddressOwner,
			Amount:         amountBig,
			CoinType:       CoinType(balChange.CoinType),
			Time:           txTimestamp,
		}
		txList = append(txList, tx)
	}
	return txList, nil
}

func (tb *TransactionBlock) Status() Status {
	return tb.Effects.Status

}

func (tb *TransactionBlock) GetGasUsed() (*big.Int, error) {
	gasUsed := tb.Effects.GasUsed
	totalBig := big.NewInt(0)
	computationCost, ok := big.NewInt(0).SetString(gasUsed.ComputationCost, 10)
	if !ok {
		return totalBig, fmt.Errorf("gasused parse big error: %v", gasUsed)
	}

	totalBig = big.NewInt(0).Add(totalBig, computationCost)
	storageCost, ok := big.NewInt(0).SetString(gasUsed.StorageCost, 10)
	if !ok {
		return totalBig, fmt.Errorf("gasused parse big error: %v", gasUsed)
	}
	totalBig = big.NewInt(0).Add(totalBig, storageCost)

	storageRebate, ok := big.NewInt(0).SetString(gasUsed.StorageRebate, 10)
	if !ok {
		return big.NewInt(0), fmt.Errorf("gasused parse big error: %v", gasUsed)
	}
	totalBig = big.NewInt(0).Sub(totalBig, storageRebate)

	//nonRefundableStorageFee, ok := big.NewInt(0).SetString(gasUsed.NonRefundableStorageFee, 10)
	//if !ok {
	//	return totalBig, fmt.Errorf("gasused parse big error: %v", gasUsed)
	//}
	//totalBig = big.NewInt(0).Add(totalBig, nonRefundableStorageFee)
	return totalBig, nil

}

type Inputs struct {
	Type      string      `json:"type"`
	ValueType string      `json:"valueType"`
	Value     interface{} `json:"value"`
}

type Transactions struct {
	SplitCoins      []interface{} `json:"SplitCoins,omitempty"`
	TransferObjects []interface{} `json:"TransferObjects,omitempty"`
}

type DataTransaction struct {
	Kind         string         `json:"kind"`
	Inputs       []Inputs       `json:"inputs"`
	Transactions []Transactions `json:"transactions"`
}

type Payment struct {
	ObjectID string `json:"objectId"`
	Version  int    `json:"version"`
	Digest   string `json:"digest"`
}

type GasData struct {
	Payment []Payment `json:"payment"`
	Owner   string    `json:"owner"`
	Price   string    `json:"price"`
	Budget  string    `json:"budget"`
}

type Data struct {
	//MessageVersion string          `json:"messageVersion"`
	//Transaction DataTransaction `json:"transaction"`
	Sender  string  `json:"sender"`
	GasData GasData `json:"gasData"`
}

type Transaction struct {
	Data         Data     `json:"data"`
	TxSignatures []string `json:"txSignatures"`
}

type Status struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type GasUsed struct {
	ComputationCost         string `json:"computationCost"`
	StorageCost             string `json:"storageCost"`
	StorageRebate           string `json:"storageRebate"`
	NonRefundableStorageFee string `json:"nonRefundableStorageFee"`
}

type ModifiedAtVersions struct {
	ObjectID       string `json:"objectId"`
	SequenceNumber string `json:"sequenceNumber"`
}

type Owner struct {
	AddressOwner string `json:"AddressOwner"`
}

type Reference struct {
	ObjectID string `json:"objectId"`
	Version  int    `json:"version"`
	Digest   string `json:"digest"`
}

type Created struct {
	Owner     Owner     `json:"owner"`
	Reference Reference `json:"reference"`
}

type Mutated struct {
	Owner     Owner     `json:"owner"`
	Reference Reference `json:"reference"`
}

type GasObject struct {
	Owner     Owner     `json:"owner"`
	Reference Reference `json:"reference"`
}

type Effects struct {
	MessageVersion string  `json:"messageVersion"`
	Status         Status  `json:"status"`
	ExecutedEpoch  string  `json:"executedEpoch"`
	GasUsed        GasUsed `json:"gasUsed"`
	//ModifiedAtVersions []ModifiedAtVersions `json:"modifiedAtVersions"`
	TransactionDigest string `json:"transactionDigest"`
	//Created            []Created            `json:"created"`
	//Mutated            []Mutated            `json:"mutated"`
	//GasObject    GasObject `json:"gasObject"`
	//Dependencies []string  `json:"dependencies"`
}

type ObjectChanges struct {
	Type            string `json:"type"`
	Sender          string `json:"sender"`
	Owner           Owner  `json:"owner"`
	ObjectType      string `json:"objectType"`
	ObjectID        string `json:"objectId"`
	Version         string `json:"version"`
	PreviousVersion string `json:"previousVersion,omitempty"`
	Digest          string `json:"digest"`
}

type BalanceChanges struct {
	Owner    Owner  `json:"owner"`
	CoinType string `json:"coinType"`
	Amount   string `json:"amount"`
}
