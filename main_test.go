package main

import (
	"bytes"
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/romankarpowich/ozon/app/router"
	"github.com/romankarpowich/ozon/models"
	"github.com/romankarpowich/ozon/repository"
	"github.com/romankarpowich/ozon/repository/db"
	"github.com/romankarpowich/ozon/repository/memory"
	"github.com/romankarpowich/ozon/service"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

func TestMain(m *testing.M) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln(err)
	}
	router.InitRouter()

	*repository.MemoryType = "postgresql"
	db.InitDb()
	models.TableName = "shorts_test"

	existTable()
	code := m.Run()
	clearTable()

	*repository.MemoryType = "in-memory"
	memory.InitMemory()
	code = m.Run()

	os.Exit(code)
}

func existTable() {
	if *repository.MemoryType == "postgresql" {
		if _, err := db.Store.Exec(tableCreationQuery); err != nil {
			log.Fatal(err)
		}
	}

}

func clearTable() {
	if *repository.MemoryType == "postgresql" {
		db.Store.Exec("DELETE FROM shorts_test;")
		db.Store.Exec("ALTER SEQUENCE shorts_test_id_seq RESTART WITH 1")
	} else {
		memory.Store.Clear()
	}
}

const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS shorts_test (
                                      id BIGSERIAL,
                                      short varchar(10),
    source varchar(255),
    CONSTRAINT shorts_test_pkey PRIMARY KEY (id)
    );

select nextval('shorts_test_id_seq');
`

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func addShorts(count int) {
	if count < 1 {
		count = 1
	}

	for i := 1; i <= count; i++ {
		s := service.UrlGenerator
		if *repository.MemoryType == "postgresql" {
			db.Store.Exec("INSERT INTO shorts_test(source, short) VALUES($1, $2)", "https://youtube.com/" + strconv.Itoa(i), s.GetShortByID(uint64(i)))
		} else {
			memory.Store.Set("https://youtube.com/" + strconv.Itoa(i), s.GetShortByID(uint64(i)))
		}
	}
}

func execRequest(request *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.Mux.ServeHTTP(rr, request)
	return rr
}

func TestEmptyTable(t *testing.T) {
	clearTable()
	req, _ := http.NewRequest("GET", "/aaaaaaaaaa", nil)
	checkResponseCode(t, http.StatusNotFound, execRequest(req).Code)
}

func TestCreateShort(t *testing.T) {

	clearTable()

	var jsonStr = []byte(`{"source":"https://youtube.com"}`)
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := execRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	m := make(map[string]interface{})
	err := json.Unmarshal(response.Body.Bytes(), &m)
	if err != nil {
		t.Errorf(err.Error())
	}
	if m["source"] != "https://youtube.com" {
		t.Errorf("Expected short source to be 'https://youtube.com'. Got '%s'", m["source"])
	}

	if m["short"] != "aaaaaaaaaa" {
		t.Errorf("Expected short short to be 'aaaaaaaaaa'. Got '%s'", m["short"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected short ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetShortByShort(t *testing.T) {
	clearTable()
	addShorts(1)
	req, _ := http.NewRequest("GET", "/aaaaaaaaaa", nil)

	checkResponseCode(t, http.StatusOK, execRequest(req).Code)
}
