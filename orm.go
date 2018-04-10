package goorm

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	upper "upper.io/db.v3"
)

type orm struct {
	alias *alias
}

// NewOrm create new orm
func NewOrm() OperateSet {
	o := new(orm)
	err := o.Using("default")
	if err != nil {
		panic(err)
	}
	return o
}

// switch to another registered database driver by given name.
func (o *orm) Using(name string) error {

	if al, ok := dataBaseCache.get(name); ok {
		o.alias = al
	} else {
		return fmt.Errorf("<Ormer.Using> unknown db alias name `%s`", name)
	}
	return nil
}

//registerModel register models.
// PrefixOrSuffix means table name prefix or suffix.
// isPrefix whether the prefix is prefix or suffix
func registerModel(PrefixOrSuffix string, model interface{}, isPrefix bool) {
	val := reflect.ValueOf(model)
	typ := reflect.Indirect(val).Type()
	if val.Kind() != reflect.Ptr {
		panic(fmt.Errorf("<orm.RegisterModel> cannot use non-ptr model struct `%s`", getFullName(typ)))
	}
	// For this case:
	// m := &Model{}
	// registerModel(&m)
	if typ.Kind() == reflect.Ptr {
		panic(fmt.Errorf("<orm.RegisterModel> only allow ptr model struct, it looks you use two reference to the struct `%s`", typ))
	}

	table := getTableName(val)

	if PrefixOrSuffix != "" {
		if isPrefix {
			table = PrefixOrSuffix + table
		} else {
			table = table + PrefixOrSuffix
		}
	}
	// models's fullname is pkgpath + struct name
	name := getFullName(typ)

	if _, ok := modelPool.getByFullName(name); ok {
		fmt.Printf("<orm.RegisterModel> model `%s` repeat register, must be unique\n", name)
		os.Exit(2)
	}
	if _, ok := modelPool.get(table); ok {
		fmt.Printf("<orm.RegisterModel> table name `%s` repeat register, must be unique\n", table)
		os.Exit(2)
	}
	mi := newModelInfo(val)
	mi.table = table
	mi.model = model
	modelPool.set(table, mi)

}

// RegisterModel register models
func RegisterModel(models ...interface{}) {
	for _, model := range models {
		registerModel("", model, true)
	}
}

// get model info and model reflect value by type
func (o *orm) getMiIndByType(typ reflect.Type) (*modelInfo, error) {
	name := getFullName(typ)
	if mi, ok := modelPool.getByFullName(name); ok {
		return mi, nil
	}
	err := fmt.Errorf("<Orm> table: `%s` not found, make sure it was registered with `RegisterModel()`", name)
	return nil, err
}

// get model info and model reflect value
func (o *orm) getMiInd(md interface{}) (*modelInfo, error) {
	val := reflect.ValueOf(md)
	ind := reflect.Indirect(val)
	typ := ind.Type()
	if val.Kind() != reflect.Ptr {
		err := fmt.Errorf("<Orm> cannot use non-ptr model struct `%s`", getFullName(typ))
		return nil, err
	}
	name := getFullName(typ)
	if mi, ok := modelPool.getByFullName(name); ok {
		return mi, nil
	}
	err := fmt.Errorf("<Orm> table: `%s` not found, make sure it was registered with `RegisterModel()`", name)
	return nil, err
}

func (o *orm) Insert(model interface{}) (int64, error) {
	mi, err := o.getMiInd(model)
	if err != nil {
		return 0, err
	}
	_id, err := (*o.alias.DB).Collection(mi.table).Insert(model)
	if err != nil {
		return 0, err
	}
	sid := fmt.Sprint(_id)
	id, _ := strconv.ParseInt(sid, 10, 64)
	return id, nil
}

func (o *orm) Update(model interface{}, cols ...string) error {
	mi, err := o.getMiInd(model)
	if err != nil {
		return err
	}
	err = (*o.alias.DB).Collection(mi.table).UpdateReturning(model)
	if err != nil {
		return err
	}
	return nil
}
func (o *orm) Delete(model interface{}) (int64, error) {
	mi, err := o.getMiInd(model)
	if err != nil {
		return 0, err
	}
	val := reflect.Indirect(reflect.ValueOf(model))
	typ := reflect.Indirect(reflect.ValueOf(model)).Type()
	n := typ.NumField()
	var id interface{}
	for i := 0; i < n; i++ {
		if vtag, ok := typ.Field(i).Tag.Lookup("db"); ok {
			if vtag == "id" {
				id = val.Field(i).Interface()
				break
			}
		}
	}
	if id == nil {
		return 0, fmt.Errorf("<Orm> table: `%s` missing primary key id field", getFullName(typ))
	}
	findR := (*o.alias.DB).Collection(mi.table).Find(upper.Cond{
		"id": id,
	})
	c, err := findR.Count()
	if err != nil {
		return 0, err
	}
	if c > 0 {
		err = findR.Delete()
		if err != nil {
			return 0, err
		}
	}
	return int64(c), nil
}

//Select
//
func (o *orm) Select(models interface{}, fields map[interface{}]interface{}) error {
	val := reflect.ValueOf(models)
	if val.Kind() != reflect.Ptr || reflect.Indirect(val).Kind() != reflect.Slice {
		err := fmt.Errorf("<Orm> select:parameter `models` must be of type slice pointer")
		return err
	}
	if reflect.Indirect(val).Type().Elem().Kind() == reflect.Ptr {
		err := fmt.Errorf("<Orm> select:parameter `models` is not allowed to be a pointer slice")
		return err
	}
	mi, err := o.getMiIndByType(reflect.Indirect(val).Type().Elem())
	if err != nil {
		return err
	}
	sele := upper.Cond{}
	if fields != nil && len(fields) > 0 {
		for k, v := range fields {
			sele[k] = v
		}
	}
	err = (*o.alias.DB).Collection(mi.table).Find(sele).All(models)
	if err != nil {
		return err
	}

	return nil
}
func (o *orm) SelectLimit(models interface{}, fields map[interface{}]interface{}, offset int64, limit int64) error {
	val := reflect.ValueOf(models)
	if val.Kind() != reflect.Ptr || reflect.Indirect(val).Kind() != reflect.Slice {
		err := fmt.Errorf("<Orm> select:parameter `models` must be of type slice pointer")
		return err
	}
	if reflect.Indirect(val).Type().Elem().Kind() == reflect.Ptr {
		err := fmt.Errorf("<Orm> select:parameter `models` is not allowed to be a pointer slice")
		return err
	}
	mi, err := o.getMiIndByType(reflect.Indirect(val).Type().Elem())
	if err != nil {
		return err
	}
	sele := upper.Cond{}
	if fields != nil && len(fields) > 0 {
		for k, v := range fields {
			sele[k] = v
		}
	}
	err = (*o.alias.DB).Collection(mi.table).Find(sele).Offset(int(offset)).Limit(int(limit)).All(models)
	if err != nil {
		return err
	}

	return nil
}
func (o *orm) One(model interface{}, fields map[interface{}]interface{}) error {
	mi, err := o.getMiInd(model)
	if err != nil {
		return err
	}
	sele := upper.Cond{}
	if fields != nil && len(fields) > 0 {
		for k, v := range fields {
			sele[k] = v
		}
	}
	err = (*o.alias.DB).Collection(mi.table).Find(sele).One(model)
	if err != nil {
		if err == upper.ErrNoMoreRows {
			return ErrNoMoreRows
		}
		return err
	}
	return nil
}
func (o *orm) Count(model interface{}, fields map[interface{}]interface{}) (int64, error) {
	mi, err := o.getMiInd(model)
	if err != nil {
		return 0, err
	}
	sele := upper.Cond{}
	if fields != nil && len(fields) > 0 {
		for k, v := range fields {
			sele[k] = v
		}
	}
	//sliceOfStructs
	c, err := (*o.alias.DB).Collection(mi.table).Find(sele).Count()
	if err != nil {
		return 0, err
	}

	return int64(c), nil
}
