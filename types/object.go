package types

type ObjectInfo struct {
	ObjectID            string `json:"objectId"`
	Version             int    `json:"version"`
	Digest              string `json:"digest"`
	Type                string `json:"type"`
	Owner               Owner  `json:"owner"`
	PreviousTransaction string `json:"previousTransaction"`
}
