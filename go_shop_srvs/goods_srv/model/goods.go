package model

//类型， 这个字段是否能为null， 这个字段应该设置为可以为null还是设置为空， 0
//实际开发过程中 尽量设置为不为null
//https://zhuanlan.zhihu.com/p/73997266
//这些类型使用int32还是int
type Category struct {
	BaseModel
	Name             string      `gorm:"type:varchar(20);not null" json:"name"`
	ParentCategoryID int32       `json:"parent"`
	ParentCategory   *Category   `json:"-"`
	SubCategory      []*Category `gorm:"foreignKey:ParentCategoryID;references:ID" json:"sub_category"`
	Level            int32       `gorm:"type:int;not null;default:1" json:"level"`
	IsTab            bool        `gorm:"default:false;not null" json:"is_tab"`
}

type Brands struct {
	BaseModel
	Name  string `gorm:"type:varchar(20);not null"`
	Logo  string `gorm:"type:varchar(200);default:'';not null"`
}

type GoodsCategoryBrand struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Category Category

	BrandsID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Brands Brands
}

func (GoodsCategoryBrand) TableName() string {
	return "goodscategorybrand"	
}

type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null"`
	Url string `gorm:"type:varchar(200);not null"`
	Index int32 `gorm:"type:int;default:1;not null"`
}


type Goods struct {
	BaseModel

	CategoryID int32 `gorm:"type:int;not null"`
	Category Category
	BrandsID int32 `gorm:"type:int;not null"`
	Brands Brands

	OnSale bool `gorm:"default:false;not null"`
	ShipFree bool `gorm:"default:false;not null"`
	IsNew bool `gorm:"default:false;not null"`
	IsHot bool `gorm:"default:false;not null"`

	Name  string `gorm:"type:varchar(50);not null"`
	GoodsSn string `gorm:"type:varchar(50);not null"` // 商家内部编号,可以去库存中拿
	ClickNum int32 `gorm:"type:int;default:0;not null"`
	SoldNum int32 `gorm:"type:int;default:0;not null"` 
	FavNum int32 `gorm:"type:int;default:0;not null"` // 收藏数
	MarketPrice float32 `gorm:"not null"`   // 市场价
	ShopPrice float32 `gorm:"not null"`		// 售卖价
	GoodsBrief string `gorm:"type:varchar(100);not null"` // 简介
	Images GormList `gorm:"type:varchar(1000);not null"`  // []string 1.可以自定义类型；2.可以建一张表,查询时需要join,数据量大时会有很大性能降低
	DescImages GormList `gorm:"type:varchar(1000);not null"`	// 应该会转化成text,varchar长度为256
	GoodsFrontImage string `gorm:"type:varchar(200);not null"`
}
