package main

import (
	"fmt"
	"github.com/boltdb/bolt"
)

func main() {
	//打开数据库
	//句柄，数据库对象handler
	//func Open(path string, mode os.FileMode, options *Options) (*DB, error) {
	db, err := bolt.Open("test.db", 0600, nil /*超时相关配置*/)
	if err != nil {
		fmt.Println("bolt.Open err:", err)
		return
	}

	defer db.Close()

	//打开桶，Bucket, 如果没有桶，则需要手动创建
	db.Update(func(tx *bolt.Tx) error {

		//尝试打开一个bucket
		b := tx.Bucket([]byte("blockBucket"))

		if b == nil {
			//当前的bucket不存在，需要创建
			b, err = tx.CreateBucket([]byte("blockBucket"))
			if err != nil {
				fmt.Println(" tx.CreateBucket err:", err)
				return err
			}
		}

		//此时b一定是非nil的
		//添加
		_ = b.Put([]byte("key1"), []byte("hello"))
		_ = b.Put([]byte("key2"), []byte("world"))
		_ = b.Put([]byte("key2"), []byte("WORLD")) //覆盖

		//读取
		v1 := b.Get([]byte("key1"))
		v2 := b.Get([]byte("key2"))
		v3 := b.Get([]byte("key3")) //返回空值

		fmt.Println("v1:", string(v1))
		fmt.Println("v2:", string(v2))
		fmt.Println("v3:", string(v3))

		return nil
	})
}
