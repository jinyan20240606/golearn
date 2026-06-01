package main

import (
	"sync"
	"sync/atomic"
)

type DBPool struct {
	Host     string
	Port     int
	UserName string
}

var dbPoolIns *DBPool
var lock sync.Mutex
var initialed uint32

// 有问题的方法，并发
// 这种单例写法有很大的并发安全问题：同时请求时，会创建多个实例-------- 并发安全问题！！！
// func GetDBPoolOld() any {
// 	if dbPoolIns == nil {
// 		dbPoolIns = &DBPool{}
// 	}
// 	return dbPoolIns
// }

// 这种普通加锁的方式： 功能上没有问题，但是对于一些频繁调用时就性能不好，因为加了之后它相当于串行
// func GetDBPool() any {
// 	lock.Lock()
// 	defer lock.Unlock()
// 	if dbPoolIns == nil {
// 		dbPoolIns = &DBPool{}
// 	}
// 	return dbPoolIns
// }

//把锁加在条件判断里，解决性能问题，就第一次需要排队，后续都直接返回，----- 但是高并发下，这种写法有bug
// func GetDBPool() any {

// 	if dbPoolIns == nil {
// 	    lock.Lock()
// 	    defer lock.Unlock()
// 		dbPoolIns = &DBPool{}
// 	}
// 	return dbPoolIns
// }

// goroutine1 进来， 实例化dbPoolIns = &DBPool{} 进行到一半，goroutine2 进来，dbPoolIns读到dbPoolIns != nil，返回dbPoolIns
// 完善解决并发安全的写法1：
func GetDBPool() *DBPool {
	if atomic.LoadUint32(&initialed) == 1 {
		return dbPoolIns
	}
	lock.Lock()
	defer lock.Unlock()

	if initialed == 0 {
		dbPoolIns = &DBPool{}
		atomic.StoreUint32(&initialed, 1)
	}

	return dbPoolIns
}

// 更好写法2: go专门提供了创建单例模式的方法，原理几乎与上面麻烦的写法一致
var once sync.Once

func GetDBPool2() *DBPool {
	once.Do(func() {
		dbPoolIns = &DBPool{}
	})
	return dbPoolIns
}
