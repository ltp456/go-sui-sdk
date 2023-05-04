package types

type Object struct {
	ObjectID string `json:"objectId"`
	Version  string `json:"version"`
	Digest   string `json:"digest"`
}
type ObjData struct {
	Data Object `json:"data"`
}

type ObjectInfo struct {
	Data        []ObjData `json:"data"`
	NextCursor  string    `json:"nextCursor"`
	HasNextPage bool      `json:"hasNextPage"`
}

type CoinData struct {
	CoinType            string      `json:"coinType"`
	CoinObjectID        string      `json:"coinObjectId"`
	Version             string      `json:"version"`
	Digest              string      `json:"digest"`
	Balance             string      `json:"balance"`
	LockedUntilEpoch    interface{} `json:"lockedUntilEpoch"`
	PreviousTransaction string      `json:"previousTransaction"`
}

type CoinObj struct {
	Data        []CoinData `json:"data"`
	NextCursor  string     `json:"nextCursor"`
	HasNextPage bool       `json:"hasNextPage"`
}
