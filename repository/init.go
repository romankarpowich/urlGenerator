package repository

import (
	"github.com/romankarpowich/ozon/repository/db"
	"github.com/romankarpowich/ozon/repository/memory"
)

var (
	MemoryType *string = new(string)
)

func InitStore() {
	switch *MemoryType {
	case "in-memory":
		memory.InitMemory()
		return
	case "postgresql":
		fallthrough
	default:
		db.InitDb()
		return
	}
}
