package entities

import "time"

const TableNameProduct = "products"

// ProductEntity is the persistence model for the products table.
type ProductEntity struct {
	ID          uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Name        string    `gorm:"column:name;size:255;not null" json:"name"`
	Description string    `gorm:"column:description;type:text" json:"description"`
	Price       float64   `gorm:"column:price;not null" json:"price"`
	Stock       int       `gorm:"column:stock;not null;default:0" json:"stock"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for GORM.
func (*ProductEntity) TableName() string {
	return TableNameProduct
}
