package goorm

import(
	"context"
)

type IOrm interface {
	IOperateSet
	ITransaction
}

type ITransaction interface {
	Tx(ctx context.Context,fn func(tx ITx) error) error
}

type ITx interface {
	//goorm.OperateSet adds operation methods to the transaction.
	IOperateSet
	//Rollbacker add Rollback method to the transaction.
	Rollbacker
	//Commiter add Commiter method to the transaction.
	Commiter
}

type Rollbacker interface {
	// Rollback discards all the instructions on the current transaction.
	Rollback() error
}

type Commiter interface {
	// Commit commits the current transaction.
	Commit() error
}
type IOperateSet interface {
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
