package models

import (
	"github.com/goravel/framework/database/orm"
)

type Blog struct {
	orm.Model
	Title   string
	Content string
	UserId  int64
	Type    int64
}

func (r *Blog) TableName() string {
	return "blogs"
}
func (r *Blog) Connection() string {
	return "mysql"
}
