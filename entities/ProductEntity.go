package entities

import "time"

// TableNameProduct is the database table used for product records.
const TableNameProduct = "product"

type ProductEntity struct {
	ID          uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Name        string    `gorm:"column:name;size:255;not null" json:"name"`
	Description string    `gorm:"column:description;type:text" json:"description"`
	Price       float64   `gorm:"column:price;not null" json:"price"`
	Category    string    `gorm:"column:category;size:50" json:"category"`
	ImageUrl    string    `gorm:"column:imageUrl;type:text" json:"imageUrl"`
	Stock       int       `gorm:"column:stock;not null;default:0" json:"stock"`
	Rating      float64   `gorm:"column:rating;type:decimal(3,2)" json:"rating"`
	ReviewCount int       `gorm:"column:reviewCount;default:0" json:"reviewCount"`
	SKU         *string   `gorm:"column:sku;unique;size:50" json:"sku"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (*ProductEntity) TableName() string {
	return TableNameProduct
}
