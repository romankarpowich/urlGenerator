package memory

import (
	_ "github.com/lib/pq"
	"sync"
)

var (
	Store *memory
)

type getterSetter interface {
	Get(short string) (string, bool)
	Set(source, short string)
	GetId() uint64
}

type memory struct {
	mut  sync.Mutex
	data map[string]string
	l    uint64
	getterSetter
}

func InitMemory() {
	Store = new(memory)
	Store.data = make(map[string]string)
	Store.l = 0
}

func (m *memory) Get(short string) (string, bool) {
	item, ok := m.data[short]
	return item, ok
}

func (m *memory) Set(source, short string) {
	m.mut.Lock()
	m.data[short] = source
	m.l++
	m.mut.Unlock()
}

func (m *memory) GetId() uint64 {
	return m.l
}

func (m *memory) Clear() {
	m.mut.Lock()
	m.data = make(map[string]string)
	m.l = 0
	m.mut.Unlock()
}
