package dao

import (
	"alice-chatgpt/conversation"
	"fmt"
	"github.com/dgraph-io/badger/v3"
)

type Dao interface {
	InitDatabase() error
	Save(storage *conversation.CStorage) error
	Search(id string) *conversation.CStorage
	Close() error
}
type BadgerDao struct {
	db *badger.DB
}

func (dao *BadgerDao) InitDatabase() error {
	db, err := badger.Open(badger.DefaultOptions("/data"))
	if err != nil {
		return err
	}
	dao.db = db
	return nil
}

func (dao *BadgerDao) Save(storage *conversation.CStorage) error {
	return dao.db.Update(func(txn *badger.Txn) error {
		bytes, err := storage.ToJsonBytes()
		if err != nil {
			return err
		}
		err = txn.Set([]byte(storage.Id), bytes)
		return err
	})
}

func (dao *BadgerDao) Search(id string) *conversation.CStorage {
	var jsonBytes []byte
	err := dao.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(id))
		if err != nil {
			return err
		}
		jsonBytes, err = item.ValueCopy(nil)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		return nil
	}
	cStorage, err2 := conversation.FromJsonBytes(jsonBytes)
	if err2 != nil {
		fmt.Println(err2.Error())
		return nil
	}
	return cStorage
}

func (dao *BadgerDao) Close() error {
	return dao.db.Close()
}
