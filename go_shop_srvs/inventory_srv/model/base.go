package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type GormList []string

func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}

type BaseModel struct {
	ID        int32     `gorm:"primary_key";type:int`
	CreatedAt time.Time `gorm:"column:add_time" json:"-"`
	UpdatedAt time.Time `gorm:"column:update_time" json:"-"`
	// DeleteAt  gorm.DeletedAt `json:"-"`   老查询报不存在
	IsDeleted bool `json:"-"`
}
