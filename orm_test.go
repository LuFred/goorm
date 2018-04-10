package goorm

import (
	"testing"
)

var o OperateSet

func init() {
	d := &DNA{}
	RegisterModel(d)
	RegisterDataBase("default", "mysql", "host", "database", "user", "pwd")
	o = NewOrm()
	o.Using("default")
}
func TestInsert(t *testing.T) {
	id, err := o.Insert(&DNA{
		UserID:      12,
		OriginDNAID: 2,
		StyleName:   "xxx",
		Name:        "x",
		Thumbnail:   "x",
		StyleID:     1,
		Introduce:   "xx",
		Price:       12,
		GMTCreate:   123,
	})
	t.Errorf("id %d err %v", id, err)
}

func TestUpdate(t *testing.T) {
	err := o.Update(&DNA{
		ID:          12,
		UserID:      123333,
		OriginDNAID: 2,
		StyleName:   "xxx",
		Name:        "x",
		Thumbnail:   "x",
		StyleID:     1,
		Introduce:   "xx",
		Price:       12,
		GMTCreate:   123,
	})
	t.Errorf("err %v", err)
}
func TestDelete(t *testing.T) {
	r, err := o.Delete(&DNA{
		ID:          2,
		UserID:      123333,
		OriginDNAID: 2,
		StyleName:   "xxx",
		Name:        "x",
		Thumbnail:   "x",
		StyleID:     1,
		Introduce:   "xx",
		Price:       12,
		GMTCreate:   123,
	})
	t.Errorf("err %d,%v", r, err)
}

func TestSelect(t *testing.T) {
	var ds []DNA
	err := o.Select(&ds, map[interface{}]interface{}{
		"id": 14,
	})
	t.Errorf("---%v", ds)

	t.Errorf("id %d err %v", len(ds), err)
}

func TestSelectLimit(t *testing.T) {

	var ds []DNA
	err := o.SelectLimit(&ds, map[interface{}]interface{}{}, 0, 2)
	t.Errorf("---%v", ds)

	t.Errorf("id %d err %v", len(ds), err)
}

func TestCount(t *testing.T) {
	d := &DNA{}

	count, err := o.Count(d, map[interface{}]interface{}{
		"name": "x1",
	})

	t.Errorf("c %d err %v", count, err)
}

func TestOne(t *testing.T) {
	d := &DNA{}

	err := o.One(d, map[interface{}]interface{}{
		"name": "x",
	})
	t.Errorf("---%v", d)

	t.Errorf("id %d err %v", d, err)
}
