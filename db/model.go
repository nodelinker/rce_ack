package db

import "time"

// PwnDataCache [...]
type PwnDataCache struct {
	ID          int    `gorm:"primaryKey;autoIncrement;column:id" json:"-"`
	ClientIP    string `gorm:"column:client_ip" json:"clientIp"`
	ClientPort  int    `gorm:"column:client_port" json:"clientPort"`
	RawDataHash string `gorm:"column:raw_data_hash" json:"rawDataHash"`
	RawData     string `grom:"column:raw_data" json:"rawData"`

	Created   int64     `gorm:"autoUpdateTime:milli"` // Use unix milli seconds as updating time
	CreatedAt time.Time `gorm:"column:create_at" json:"createAt"`
	UpdatedAt time.Time `gorm:"column:update_at" json:"updateAt"`
}

// TableName get sql table name.获取数据库表名
func (m *PwnDataCache) TableName() string {
	return "pwn_data_cache"
}
