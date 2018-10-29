package goorm

type OperateSet interface {
	Insert(model interface{}) (int64, error)
	Update(model interface{}, cols ...string) error
	Delete(model interface{}) (int64, error)
	Select(models interface{}, fields Cond, orderby ...interface{}) error
	SelectLimit(models interface{}, fields Cond, offset int64, limit int64, orderby ...interface{}) error
	One(model interface{}, fields Cond) error
	Count(model interface{}, fields Cond) (int64, error)
	Exec(sql string, args ...interface{}) (int64, error)
	Using(name string) error
	QuerySQL(dests interface{}, sql string, args ...interface{}) error
	ExecSQL(dest interface{}, sql string, args ...interface{}) error
}
