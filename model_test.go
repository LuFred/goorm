package goorm

import (
	"reflect"
	"testing"
)

type DNA struct {
	ID          int32  `db:"id"`            //自增主键id
	UserID      int64  `db:"user_id"`       //用户id
	OriginDNAID int32  `db:"origin_dna_id"` //原始dna id
	StyleName   string `db:"style_name"`    //名称
	Name        string `db:"name"`          //名称
	Thumbnail   string `db:"thumbnail"`     //缩略图
	StyleID     int32  `db:"style_id"`      //类型id
	Introduce   string `db:"introduce"`     //介绍
	Price       int64  `db:"price"`         //价格
	GMTCreate   int64  `db:"gmt_create"`    //创建时间
	GMTModified int64  `db:"gmt_modified"`  //修改时间
}

func (d *DNA) TableName() string {
	return "cm_dna"
}

func TestRegisterModel(t *testing.T) {
	d := DNA{}
	val := reflect.ValueOf(d)
	reflect.Indirect(val)
	val.Elem()
	t.Error(val.Elem() == reflect.Indirect(val))
	// t.Error("--" + reflect.Indirect(val).Type().PkgPath())
	// RegisterModel(&d)
}
