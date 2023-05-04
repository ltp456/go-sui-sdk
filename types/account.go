package types

type AllBalance struct {
	CoinType        string        `json:"coinType"`
	CoinObjectCount int           `json:"coinObjectCount"`
	TotalBalance    string        `json:"totalBalance"`
	LockedBalance   LockedBalance `json:"lockedBalance"`
}

type LockedBalance struct {
}
