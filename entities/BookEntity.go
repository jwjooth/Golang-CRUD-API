package entities

import "time"

const TableNameBook = "books"

type BookEntity struct {
	ID         uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CategoryId *uint     `gorm:"column:category_id" json:"category_id"`
	Title      string    `gorm:"column:title;size:255;not null" json:"title"`
	Author     string    `gorm:"column:author;size:255;not null" json:"author"`
	Stock      int       `gorm:"column:stock;not null;default:0" json:"stock"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (*BookEntity) TableName() string {
	return TableNameBook
}
