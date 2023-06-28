package model

type Router struct {
	ID              int64        `gorm:"primary_key;not_null;auto_increment"`
	RouterName      string       `json:"router_name"`
	RouterNamespace string       `json:"router_namespace"`
	RouterHost      string       `json:"router_host"`
	RouterPath      []RouterPath `gorm:"ForeignKey:RouterID" json:"router_path"`
}
