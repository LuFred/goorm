package goorm

type OperateSet interface {
	Insert(model interface{}) (int64, error)
	Update(model interface{}, cols ...string) error
	Delete(model interface{}) (int64, error)
	Select(models interface{}, fields map[interface{}]interface{}) error
	SelectLimit(models interface{}, fields map[interface{}]interface{}, offset int64, limit int64) error
	One(model interface{}, fields map[interface{}]interface{}) error
	Count(model interface{}, fields map[interface{}]interface{}) (int64, error)
	Using(name string) error
}
