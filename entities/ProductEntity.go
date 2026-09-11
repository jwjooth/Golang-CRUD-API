package entities

import "time"

const TableNameProduct = "products"

type ProductEntity struct {
	Id          int       `gorm:"primaryKey;column:id" json:"id"`
	Name        string    `gorm:"column:name" json:"name"`
	Description string    `gorm:"column:description" json:"description"`
	Price       int64     `gorm:"column:price" json:"price"`
	Stock       int       `gorm:"column:stock" json:"stock"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (*ProductEntity) TableName() string {
	return TableNameProduct
}
