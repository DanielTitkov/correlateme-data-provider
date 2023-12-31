// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameMetric = "metrics"

// Metric mapped from table <metrics>
type Metric struct {
	ID        string    `gorm:"column:id;primaryKey;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"column:name;not null" json:"name"`
	Desc      string    `gorm:"column:desc" json:"desc"`
	IsPublic  bool      `gorm:"column:is_public;not null" json:"is_public"`
	IsSystem  bool      `gorm:"column:is_system;not null" json:"is_system"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	Meta      string    `gorm:"column:meta" json:"meta"`
	IsActive  bool      `gorm:"column:is_active;not null;default:true" json:"is_active"`
	UserID    string    `gorm:"column:user_id;not null" json:"user_id"`
	Kind      string    `gorm:"column:kind;not null;default:numeric" json:"kind"`
}

// TableName Metric's table name
func (*Metric) TableName() string {
	return TableNameMetric
}
