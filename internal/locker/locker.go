package locker

import (
	"log"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
)

type Element struct {
	Key   string
	Value string
}

type Locker struct {
	Key      string
	Locked   bool
	Elements []Element
	Db       *leveldb.DB
}

func (l *Locker) Connect() {
	db, err := leveldb.OpenFile("tmp/test.db", nil)
	if err != nil {
		log.Fatal("Yikes!")
	}
	l.Db = db
}

func (l *Locker) Disconnect() {
	l.Db.Close()
}

func (l *Locker) Lock() {
	l.Locked = true
}

func (l *Locker) Unlock() {
	// Require main password
	l.Locked = false
}

func (l *Locker) AddElement(key string, value string) {
	if l.Locked {
		return
	}
	l.Elements = append(l.Elements, Element{Key: key, Value: value})
	err := l.Db.Put([]byte(key), []byte(value), nil)
	if err != nil {
		log.Fatal("Yikes!")
	}
}

func (l *Locker) GetElement(key string) string {
	if l.Locked {
		return ""
	}
	for _, element := range l.Elements {
		if element.Key == key {
			foundElement := element.Value
			log.Print(foundElement)
		}
	}
	data, err := l.Db.Get([]byte(key), nil)
	if err != nil {
		log.Fatal("Yikes!")
	}
	return string(data)
}

func (l *Locker) RemoveElement(key string) {
	if l.Locked {
		return
	}
	// Find the index of the object to remove
	indexToRemove := -1
	for i, element := range l.Elements {
		if element.Key == key {
			indexToRemove = i
			break
		}
	}
	if indexToRemove != -1 {
		l.Elements = append(l.Elements[:indexToRemove], l.Elements[indexToRemove+1:]...)
	}
}

func (l *Locker) GetAllElements() []Element {
	iter := l.Db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := strings.Repeat("*", len(iter.Value()))
		l.Elements = append(l.Elements, Element{Key: string(key), Value: string(value)})
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return []Element{}
	}
	return l.Elements
}
