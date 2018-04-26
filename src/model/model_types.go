package model

type User struct {
	Id   int64  `sql:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type Relationship struct {
	Id     int64  `json:"id"`
	Uid    int64  `json:"uid"`
	Oid    int64  `json:"oid"`
	Status string `json:"status"`
	Type   string `json:"type"`
}

const (
	USER_TYPE = iota
	RELATIONSHIPS_TYPE
)

const (
	LIKED_ = iota
	DISLIKED_
	MATCHED_
)

const LIKED = "0"
const DISLIKED = "1"
const MATCHED = "2"

var StateMap map[string]int = map[string]int{
	"liked":    LIKED_,
	"disliked": DISLIKED_,
}
