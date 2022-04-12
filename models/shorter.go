package models

import (
	"errors"
	"fmt"
	"github.com/romankarpowich/ozon/repository"
	"github.com/romankarpowich/ozon/repository/db"
	"github.com/romankarpowich/ozon/repository/memory"
	"github.com/romankarpowich/ozon/service"
	"strings"
)

var (
	TableName = "shorts"
)

type shorter interface {
	Generate() error
	Get() error
}

type Short struct {
	ID     uint64 `json:"id"`
	Source string `json:"source"`
	Short  string `json:"short"`
	shorter
}

func (u *Short) Generate() error {
	id := uint64(1)

	builder := new(strings.Builder)

	defer builder.Reset()

	if *repository.MemoryType == "postgresql" {
		builder.WriteString("SELECT currval(pg_get_serial_sequence('")
		builder.WriteString(TableName)
		builder.WriteString("','id'))")

		err := db.Store.QueryRow(builder.String()).Scan(&id) // TODO разобраться почему бзе nextval (в консоли) е работает
		if err != nil {
			id = 1
		}
		builder.Reset()

		u.Short = service.UrlGenerator.GetShortByID(id)

		builder.WriteString("INSERT INTO ")
		builder.WriteString(TableName)
		builder.WriteString("(source, short) VALUES($1, $2) RETURNING id")

		err = db.Store.QueryRow(
			builder.String(),
			u.Source, u.Short).Scan(&u.ID)

		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	} else {
		u.Short = service.UrlGenerator.GetShortByID(memory.Store.GetId())
		memory.Store.Set(u.Source, u.Short)
		u.ID = memory.Store.GetId()
		return nil
	}
}

func (u *Short) Get() error {

	if *repository.MemoryType == "postgresql" {
		builder := new(strings.Builder)

		defer builder.Reset()

		builder.WriteString("SELECT id, source, short FROM ")
		builder.WriteString(TableName)
		builder.WriteString(" WHERE short=$1")

		return db.Store.QueryRow(builder.String(),
			u.Short).Scan(&u.ID, &u.Source, &u.Short)
	} else {
		item, ok := memory.Store.Get(u.Short)
		if ok {
			u.Source = item
			return nil
		}
		return errors.New("Not found")
	}
}
