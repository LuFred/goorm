# goorm
一个简单的数据库crud操作库  
对[upper](https://upper.io/) 这数据访问框架的二次封装



## 用法
```
	d := &DNA{}
	RegisterModel(d)//注册表结构
	RegisterDataBase("default", "mysql", "host", "database", "user", "pwd")//注册数据库连接对象
	o = NewOrm()
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
 
```
具体crud可看orm_test.go
