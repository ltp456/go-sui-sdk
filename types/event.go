package types

type TxID struct {
	TxDigest string `json:"txDigest"`
	EventSeq string `json:"eventSeq"`
}

type DolaUserAddress struct {
	DolaAddress []int `json:"dola_address"`
	DolaChainID int   `json:"dola_chain_id"`
}

type ParsedJSON struct {
	CallType        int             `json:"call_type"`
	Nonce           string          `json:"nonce"`
	Sender          string          `json:"sender"`
	UserAddress     []int           `json:"user_address"`
	UserChainID     int             `json:"user_chain_id"`
	DolaUserAddress DolaUserAddress `json:"dola_user_address"`
	DolaUserID      string          `json:"dola_user_id"`
}

type TxEvent struct {
	ID                TxID       `json:"id"`
	PackageID         string     `json:"packageId"`
	TransactionModule string     `json:"transactionModule"`
	Sender            string     `json:"sender"`
	Type              string     `json:"type"`
	ParsedJSON        ParsedJSON `json:"parsedJson,omitempty"`
	Bcs               string     `json:"bcs"`
}
