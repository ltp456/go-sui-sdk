package types

import "math/big"

type Tx struct {
	Sender    string   `json:"from"`
	Recipient string   `json:"to"`
	Amount    *big.Int `json:"amount"`
	Gas       *big.Int `json:"gas"`
	CoinType  CoinType `json:"coin_type"`
	Height    uint64   `json:"height"`
	Hash      string   `json:"hash"`
	TxMethod  string   `json:"tx_method"`
}

type Transaction struct {
	Certificate             Certificate `json:"certificate"`
	Effects                 Effects     `json:"effects"`
	TimestampMs             interface{} `json:"timestamp_ms"`
	ParsedData              interface{} `json:"parsed_data"`
	ConfirmedLocalExecution bool        `json:"confirmed_local_execution"`
}

func (t Transaction) Parse(index uint64) ([]Tx, error) {
	var txes []Tx
	if !t.Success() {
		return txes, nil
	}
	gasUsedBig := t.GasUsed()
	txHash := t.Effects.TransactionDigest
	sender := t.Certificate.Data.Sender

	events := t.Effects.Events
	var eventsByChangeType []Event
	for _, event := range events {
		if !event.ValidModule() {
			continue
		}
		if event.CoinBalanceChange.ChangeType == ReceiveChangeType && event.CoinBalanceChange.PackageID == SuiPackageId {
			eventsByChangeType = append(eventsByChangeType, event)
		}
	}
	if len(eventsByChangeType) == 0 {
		return txes, nil
	}

	// todo
	eventAddressMap := EventFilterByAddress(eventsByChangeType)

	for recipient, eventList := range eventAddressMap {

		eventCoinTypeMap := EventFilterByCoinType(*eventList)

		for coinType, list := range eventCoinTypeMap {

			amount := big.NewInt(0)
			for _, item := range *list {
				amount = big.NewInt(0).Add(amount, item.CoinBalanceChange.Amount)
			}
			tx := Tx{
				Sender:    sender,
				Recipient: recipient,
				Amount:    amount,
				Hash:      txHash,
				CoinType:  coinType,
				Gas:       gasUsedBig,
				Height:    index,
			}
			txes = append(txes, tx)
		}

	}
	return txes, nil

}

func EventFilterByAddress(eventList []Event) map[string]*[]Event {
	addressMap := make(map[string]*[]Event)
	for _, event := range eventList {
		if list, ok := addressMap[event.CoinBalanceChange.Owner.AddressOwner]; ok {
			*list = append(*list, event)
		} else {
			events := new([]Event)
			*events = append(*events, event)
			addressMap[event.CoinBalanceChange.Owner.AddressOwner] = events
		}
	}
	return addressMap
}

func EventFilterByCoinType(eventList []Event) map[CoinType]*[]Event {
	coinTypeMap := make(map[CoinType]*[]Event)
	for _, event := range eventList {
		if list, ok := coinTypeMap[event.CoinBalanceChange.CoinType]; ok {
			*list = append(*list, event)
		} else {
			events := new([]Event)
			*events = append(*events, event)
			coinTypeMap[event.CoinBalanceChange.CoinType] = events
		}
	}
	return coinTypeMap

}

func (e Event) ValidModule() bool {
	// todo
	module := e.CoinBalanceChange.TransactionModule
	if module == PayModule || module == PaySuiModule || module == PayAllSuiModule || module == TransferSuiModule || module == TransferObjectModule {
		return true
	}
	return false
}

func (t Transaction) GasUsed() *big.Int {
	//tmpBig := big.NewInt(0).Add(t.Effects.GasUsed.ComputationCost, t.Effects.GasUsed.StorageRebate)
	//gasUsed := big.NewInt(0).Add(tmpBig, t.Effects.GasUsed.StorageCost)
	return t.Effects.GasUsed.ComputationCost
}

func (t Transaction) Success() bool {
	return t.Effects.Status.Status == TxSuccess
}

type ObjectRef struct {
	ObjectID string `json:"objectId"`
	Version  uint64 `json:"version"`
	Digest   string `json:"digest"`
}

type Pay struct {
	Coins      []Coin     `json:"coins"`
	Recipients []string   `json:"recipients"`
	Amounts    []*big.Int `json:"amounts"`
}

type Coin struct {
	ObjectID string `json:"objectId"`
	Version  int    `json:"version"`
	Digest   string `json:"digest"`
}

type Transactions struct {
	TransferSui    TransferSui    `json:"TransferSui"`
	TransferObject TransferObject `json:"transferObject"`
	Pay            Pay            `json:"Pay"`
	PaySui         Pay            `json:"PaySui"`
	PayAllSui      Pay            `json:"PayAllSui"`
}

type TransferSui struct {
	Recipient string   `json:"recipient"`
	Amount    *big.Int `json:"amount"`
}

type GasPayment struct {
	ObjectID string `json:"objectId"`
	Version  int    `json:"version"`
	Digest   string `json:"digest"`
}

type Data struct {
	Transactions []Transactions `json:"transactions"`
	Sender       string         `json:"sender"`
	GasPayment   GasPayment     `json:"gasPayment"`
	GasBudget    uint64         `json:"gasBudget"`
}

type AuthSignInfo struct {
	Epoch      int    `json:"epoch"`
	Signature  string `json:"signature"`
	SignersMap []int  `json:"signers_map"`
}

type Certificate struct {
	TransactionDigest string       `json:"transactionDigest"`
	Data              Data         `json:"data"`
	TxSignature       string       `json:"txSignature"`
	AuthSignInfo      AuthSignInfo `json:"authSignInfo"`
}

type Status struct {
	Status TxStatus `json:"status"`
	Error  string   `json:"error"`
}

type GasUsed struct {
	ComputationCost *big.Int `json:"computationCost"`
	StorageCost     *big.Int `json:"storageCost"`
	StorageRebate   *big.Int `json:"storageRebate"`
}

type Owner struct {
	AddressOwner string `json:"AddressOwner"`
	ObjectOwner  string `json:"ObjectOwner"`
}

type Reference struct {
	ObjectID string `json:"objectId"`
	Version  int    `json:"version"`
	Digest   string `json:"digest"`
}

type Mutated struct {
	Owner     Owner     `json:"owner"`
	Reference Reference `json:"reference"`
}

type GasObject struct {
	Owner     Owner     `json:"owner"`
	Reference Reference `json:"reference"`
}

type Recipient struct {
	AddressOwner string `json:"AddressOwner"`
}

type TransferObject struct {
	PackageID         string `json:"packageId"`
	TransactionModule string `json:"transactionModule"`
	Sender            string `json:"sender"`
	//Recipient         Recipient `json:"recipient"`
	ObjectType string    `json:"objectType"`
	ObjectID   string    `json:"objectId"`
	Version    int       `json:"version"`
	Recipient  string    `json:"recipient"`
	ObjectRef  ObjectRef `json:"objectRef"`
}

type Event struct {
	TransferObject    TransferObject    `json:"transferObject"`
	CoinBalanceChange CoinBalanceChange `json:"coinBalanceChange"`
}

type Effects struct {
	Status            Status    `json:"status"`
	GasUsed           GasUsed   `json:"gasUsed"`
	TransactionDigest string    `json:"transactionDigest"`
	Mutated           []Mutated `json:"mutated"`
	GasObject         GasObject `json:"gasObject"`
	Events            []Event   `json:"events"`
}

type CoinBalanceChange struct {
	PackageID         PackageId  `json:"packageId"`
	TransactionModule Module     `json:"transactionModule"`
	Sender            string     `json:"sender"`
	ChangeType        ChangeType `json:"changeType"`
	Owner             Owner      `json:"owner"`
	CoinType          CoinType   `json:"coinType"`
	CoinObjectID      string     `json:"coinObjectId"`
	Version           int        `json:"version"`
	Amount            *big.Int   `json:"amount"`
}

// object

type ID struct {
	ID string `json:"id"`
}

type Fields struct {
	Balance *big.Int `json:"balance"`
	ID      ID       `json:"id"`
}

type ObjectData struct {
	DataType          string `json:"dataType"`
	Type              string `json:"type"`
	HasPublicTransfer bool   `json:"has_public_transfer"`
	Fields            Fields `json:"fields"`
}

type Details struct {
	Data                ObjectData `json:"data"`
	Owner               Owner      `json:"owner"`
	PreviousTransaction string     `json:"previousTransaction"`
	StorageRebate       int        `json:"storageRebate"`
	Reference           Reference  `json:"reference"`
}

type Object struct {
	Status   ObjectStatus `json:"status"`
	Details  Details      `json:"details"`
	ObjectId string       `json:"objectId"`
}

//

type UnsignedTx struct {
	TxBytes      string      `json:"txBytes"`
	InputObjects interface{} `json:"inputObjects"`
	Gas          interface{} `json:"gas"`
}

//
type SubmitTx struct {
	EffectsCert Transaction `json:"EffectsCert"`
}
