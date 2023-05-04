package types

const MinGasFee = 2000

const MaxGasFee = 50000000000

// 这在sui代码中定义
var IntentFlag = []byte{0, 0, 0}

type SignatureSchemeSerialized byte

const (
	SignatureSchemeSerializedEd25519   SignatureSchemeSerialized = 0
	SignatureSchemeSerializedSecp256k1 SignatureSchemeSerialized = 1
)

type RequestTxType string

const (
	WaitForLocalExecution RequestTxType = "WaitForLocalExecution"
	ImmediateReturn       RequestTxType = "ImmediateReturn"
	WaitForTxCert         RequestTxType = "WaitForTxCert"
	WaitForEffectsCert    RequestTxType = "WaitForEffectsCert"
)

func (rtt RequestTxType) String() string {
	return string(rtt)
}

type CoinResType string

const (
	SuiCoinResType CoinResType = "0x2::coin::Coin<0x2::sui::SUI>"
)

func (ct CoinResType) String() string {
	return string(ct)
}

type CoinType string

const (
	SuiCoinType CoinType = "0x2::sui::SUI"
)

func (ct CoinType) String() string {
	return string(ct)
}

type TxMethod string

const (
	PayAllSui TxMethod = "payAllSui"
	SuiPay    TxMethod = "suiPay"
)

func (tm TxMethod) String() string {
	return string(tm)
}

type Module string

const (
	PayModule            Module = "pay"
	PayAllSuiModule      Module = "pay_all_sui"
	PaySuiModule         Module = "pay_sui"
	TransferSuiModule    Module = "transfer_sui"
	TransferObjectModule Module = "transfer_object"
	gasModule            Module = "gas"
)

func (m Module) String() string {
	return string(m)
}

type PackageId string

const (
	SuiPackageId PackageId = "0x0000000000000000000000000000000000000002"
)

func (p PackageId) String() string {
	return string(p)
}

type ChangeType string

const (
	PayChangeType     ChangeType = "Pay"
	ReceiveChangeType ChangeType = "Receive"
	GasChangeType     ChangeType = "Gas"
)

func (ct ChangeType) String() string {
	return string(ct)
}

type ObjectStatus string

const (
	Exists ObjectStatus = "Exists"
)

func (os ObjectStatus) String() string {
	return string(os)
}

const (
	TxSuccess string = "success"
	TxFailure string = "failure"
)
