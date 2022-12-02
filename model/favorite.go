package model

import "gorm.io/gorm"

// 收藏夹
type Favorite struct {
	gorm.Model
	User      User    `gorm:"ForeignKey:UserID"` //实际场景不会用 会严重影响性能
	UserID    uint    `gorm:"not null"`
	Product   Product `gorm:"ForeignKey:ProductID"`
	ProductID uint    `gorm:"not null"`
	Boss      User    `gorm:"ForeignKey:BossID"`
	BossID    uint    `gorm:"not null"`
}
