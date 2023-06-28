package model

type RouterPath struct {
	ID                    int64  `gorm:"primary_key;not_null;auto_increment"`
	RouterID              int64  `json:"router_id"`
	RouterPathName        string `json:"router_path_name"`
	RouterBackendService  string `json:"router_backend_service"`
	RouterBackServicePort int32  `json:"router_back_service_port"`
}
