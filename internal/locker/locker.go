package locker

import (
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

type Element struct {
	key   string
	value string
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
	l.Locked = false
}

func (l *Locker) AddElement(key string, value string) {
	l.Elements = append(l.Elements, Element{key: key, value: value})
	err := l.Db.Put([]byte(key), []byte(value), nil)
	if err != nil {
		log.Fatal("Yikes!")
	}
}

func (l *Locker) GetElement(key string) string {
	for _, element := range l.Elements {
		if element.key == key {
			foundElement := element.value
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
	// Find the index of the object to remove
	indexToRemove := -1
	for i, element := range l.Elements {
		if element.key == key {
			indexToRemove = i
			break
		}
	}
	if indexToRemove != -1 {
		l.Elements = append(l.Elements[:indexToRemove], l.Elements[indexToRemove+1:]...)
	}
}
