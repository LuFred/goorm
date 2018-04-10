package goorm

import "reflect"

var (
	modelPool = &_modelPool{
		cache:           make(map[string]*modelInfo),
		cacheByFullName: make(map[string]*modelInfo),
	}
)

type _modelPool struct {
	cache           map[string]*modelInfo
	cacheByFullName map[string]*modelInfo
}

type modelInfo struct {
	pkg      string
	name     string
	fullName string
	table    string
	model    interface{}
}

func (mp *_modelPool) get(table string) (*modelInfo, bool) {
	mi, ok := mp.cache[table]
	return mi, ok
}
func (mp *_modelPool) getByFullName(name string) (*modelInfo, bool) {
	mi, ok := mp.cacheByFullName[name]
	return mi, ok
}

func (mp *_modelPool) set(table string, mi *modelInfo) {
	mp.cache[table] = mi
	mp.cacheByFullName[mi.fullName] = mi
}

func newModelInfo(val reflect.Value) (mi *modelInfo) {
	mi = &modelInfo{}
	mi.fullName = getFullName(reflect.Indirect(val).Type())
	mi.name = val.Type().Name()
	mi.pkg = val.Type().PkgPath()
	return
}
