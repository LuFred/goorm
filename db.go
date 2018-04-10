package goorm

import (
	"fmt"
	"sync"

	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/mysql"
)

var (
	dataBaseCache = &_dbCache{cache: make(map[string]*alias)}
)

// database alias cacher.
type _dbCache struct {
	mux   sync.RWMutex
	cache map[string]*alias
}

// add database alias with original name.
func (ac *_dbCache) add(name string, al *alias) (added bool) {
	ac.mux.Lock()
	defer ac.mux.Unlock()
	if _, ok := ac.cache[name]; !ok {
		ac.cache[name] = al
		added = true
	}
	return
}

// get database alias if cached.
func (ac *_dbCache) get(name string) (al *alias, ok bool) {
	ac.mux.RLock()
	defer ac.mux.RUnlock()
	al, ok = ac.cache[name]
	return
}

// get default alias.
func (ac *_dbCache) getDefault() (al *alias) {
	al, _ = ac.get("default")
	return
}

type alias struct {
	Name       string
	DriverName string
	DB         *sqlbuilder.Database
}

func RegisterDataBase(aliasName, driverName, host, database, user, pwd string) error {

	setting := &mysql.ConnectionURL{
		User:     user,
		Password: pwd,
		Host:     host,
		Database: database,
	}
	db, err := mysql.Open(setting)
	if err != nil {
		err = fmt.Errorf("register db `%s`, %s", aliasName, err.Error())
		return err
	}
	_, err = addAliasWthDB(aliasName, driverName, &db)
	if err != nil {
		err = fmt.Errorf("register db , %s", err.Error())
		return err
	}
	return nil
}

func addAliasWthDB(aliasName, driverName string, db *sqlbuilder.Database) (*alias, error) {
	al := new(alias)
	al.Name = aliasName
	al.DriverName = driverName
	al.DB = db
	err := (*al.DB).Ping()
	if err != nil {
		return nil, fmt.Errorf("register db Ping `%s`, %s", aliasName, err.Error())
	}

	if !dataBaseCache.add(aliasName, al) {
		return nil, fmt.Errorf("DataBase alias name `%s` already registered, cannot reuse", aliasName)
	}
	return al, nil
}
