package goorm

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	upper "upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

type Transaction struct{
	o *orm
	Tx sqlbuilder.Tx
}

func(t *Transaction) Using(name string) error {
	panic("transaction can not to call Using method")
}

func(t *Transaction) Insert(model interface{}) (int64, error) {
	mi, err := t.o.getMiInd(model)
	if err != nil {
		return 0, err
	}

	_id, err := t.Tx.Collection(mi.table).Insert(model)
	if err != nil {
		return 0, err
	}
	sid := fmt.Sprint(_id)
	id, _ := strconv.ParseInt(sid, 10, 64)
	return id, nil
}

func(t *Transaction) Update(model interface{}, cols ...string) error {
	mi, err := t.o.getMiInd(model)
	if err != nil {
		return err
	}
	err = t.Tx.Collection(mi.table).UpdateReturning(model)
	if err != nil {
		return err
	}
	return nil
}

func(t *Transaction)Delete(model interface{}) (int64, error) {
	mi, err := t.o.getMiInd(model)
	if err != nil {
		return 0, err
	}
	val := reflect.Indirect(reflect.ValueOf(model))
	typ := reflect.Indirect(reflect.ValueOf(model)).Type()
	n := typ.NumField()
	var id interface{}
	for i := 0; i < n; i++ {
		if vtag, ok := typ.Field(i).Tag.Lookup("db"); ok {
			if strings.ToLower(vtag) == "id" {
				id = val.Field(i).Interface()
				break
			}
		}
	}
	if id == nil {
		return 0, fmt.Errorf("<Orm> table: `%s` missing primary key id field", getFullName(typ))
	}
	findR := t.Tx.Collection(mi.table).Find(upper.Cond{
		"id": id,
	})
	defer func() {
		findR.Close()
	}()
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


func(t *Transaction) Select(models interface{}, fields Cond, orderby ...interface{}) error {
	val := reflect.ValueOf(models)
	if val.Kind() != reflect.Ptr || reflect.Indirect(val).Kind() != reflect.Slice {
		err := fmt.Errorf("<Orm> select:parameter `models` must be of type slice pointer")
		return err
	}
	if reflect.Indirect(val).Type().Elem().Kind() == reflect.Ptr {
		err := fmt.Errorf("<Orm> select:parameter `models` is not allowed to be a pointer slice")
		return err
	}
	mi, err := t.o.getMiIndByType(reflect.Indirect(val).Type().Elem())
	if err != nil {
		return err
	}
	sele := upper.Cond{}
	if fields != nil && len(fields) > 0 {
		for k, v := range fields {
			sele[k] = v
		}
	}
	if len(orderby) > 0 {
		err =t.Tx.Collection(mi.table).Find(sele).OrderBy(orderby...).All(models)
	} else {
		err = t.Tx.Collection(mi.table).Find(sele).All(models)
	}
	if err != nil {
		return err
	}

	return nil
}

func(t *Transaction) SelectLimit(models interface{}, fields Cond, offset int64, limit int64, orderby ...interface{}) error {
	if limit < 1 {
		return nil
	}
	val := reflect.ValueOf(models)
	if val.Kind() != reflect.Ptr || reflect.Indirect(val).Kind() != reflect.Slice {
		err := fmt.Errorf("<Orm> select:parameter `models` must be of type slice pointer")
		return err
	}
	if reflect.Indirect(val).Type().Elem().Kind() == reflect.Ptr {
		err := fmt.Errorf("<Orm> select:parameter `models` is not allowed to be a pointer slice")
		return err
	}
	mi, err := t.o.getMiIndByType(reflect.Indirect(val).Type().Elem())
	if err != nil {
		return err
	}
	sele := upper.Cond{}
	if fields != nil && len(fields) > 0 {
		for k, v := range fields {
			sele[k] = v
		}
	}
	if len(orderby) > 0 {
		err = t.Tx.Collection(mi.table).Find(sele).OrderBy(orderby...).Offset(int(offset)).Limit(int(limit)).All(models)
	} else {
		err =t.Tx.Collection(mi.table).Find(sele).Offset(int(offset)).Limit(int(limit)).All(models)
	}
	if err != nil {
		return err
	}

	return nil
}

func(t *Transaction) One(model interface{}, fields Cond) error {
	mi, err := t.o.getMiInd(model)
	if err != nil {
		return err
	}
	sele := upper.Cond{}
	if fields != nil && len(fields) > 0 {
		for k, v := range fields {
			sele[k] = v
		}
	}
	err = t.Tx.Collection(mi.table).Find(sele).One(model)
	if err != nil {
		if err == upper.ErrNoMoreRows {
			return ErrNoMoreRows
		}
		return err
	}
	return nil
}

func(t *Transaction) Count(model interface{}, fields Cond) (int64, error) {
	mi, err := t.o.getMiInd(model)
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
	c, err := t.Tx.Collection(mi.table).Find(sele).Count()
	if err != nil {
		return 0, err
	}

	return int64(c), nil
}


func(t *Transaction) QuerySQL(dests interface{}, sql string, args ...interface{}) error {
	val := reflect.ValueOf(dests)
	if val.Kind() != reflect.Ptr || reflect.Indirect(val).Kind() != reflect.Slice {
		err := fmt.Errorf("<Orm> select:parameter `dests` must be of type slice pointer")
		return err
	}
	if reflect.Indirect(val).Type().Elem().Kind() == reflect.Ptr {
		err := fmt.Errorf("<Orm> select:parameter `dests` is not allowed to be a pointer slice")
		return err
	}

	rows, err :=t.Tx.Query(sql, args...)
	if err != nil {
		return err
	}
	defer func() {
		rows.Close()
	}()
	iter := sqlbuilder.NewIterator(rows)
	err = iter.All(dests)
	if err != nil {
		return err
	}
	return nil
}

func(t *Transaction) Exec(sql string, args ...interface{}) (int64, error) {
	rows, err := t.Tx.Exec(sql, args...)

	if err != nil {
		return 0, err
	}
	affected, err := rows.RowsAffected()

	if err != nil {
		return 0, err
	}
	return affected, nil
}

func(t *Transaction) ExecSQL(dest interface{}, sql string, args ...interface{}) error {
	val := reflect.ValueOf(dest)
	if val.Kind() != reflect.Ptr {
		err := fmt.Errorf("<Orm> select:parameter `dest` must be of type pointer")
		return err
	}
	rows, err := t.Tx.Query(sql, args...)
	if err != nil {
		return err
	}
	defer func() {
		rows.Close()
	}()
	if rows.Next() {
		err = rows.Scan(dest)
		if err != nil {
			return err
		}
	}

	return nil
}

func(t *Transaction) Rollback() error {
	return t.Tx.Rollback()
}

func(t *Transaction) Commit() error {
	return t.Tx.Commit()
}

