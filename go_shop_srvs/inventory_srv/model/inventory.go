package model

//type Stock struct {
//	BaseModel
//	Name string
//	Address string
//}

type Inventory struct{
	BaseModel
	Goods int32 `gorm:"type:int;index"`
	Stocks int32 `gorm:"type:int"`	// 库存
	// Stock Stock  // 商品在哪个仓库,这里简化,不去考虑
	Version int32 `gorm:"type:int"` //分布式锁的乐观锁
}