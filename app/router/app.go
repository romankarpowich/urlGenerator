package router

import (
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/romankarpowich/ozon/app/response"
	"github.com/romankarpowich/ozon/models"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var (
	Mux *httprouter.Router
	re  = regexp.MustCompile("[0-9]*[a-z]*[A-Z]*_*")
)

func InitRouter() {
	Mux = httprouter.New()
	Mux.GET("/:short", show)
	Mux.POST("/", store)
}

func store(writer http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	short := new(models.Short)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(request.Body)

	bodyFromRequest, err := ioutil.ReadAll(request.Body)
	if err != nil {
		response.Failed(writer, http.StatusBadRequest, err)
		return
	}

	if err = json.Unmarshal(bodyFromRequest, &short); err != nil {
		response.Failed(writer, http.StatusBadRequest, err)
		return
	}
	if short.Source == "" {
		response.Failed(writer, http.StatusBadRequest, errors.New("Bad payload"))
		return
	}

	if err = short.Generate(); err != nil {
		response.Failed(writer, http.StatusBadRequest, err)
		return
	}
	response.Success(writer, http.StatusCreated, short)
}

func show(writer http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	val := ps.ByName("short")
	if re.MatchString(val) && len(val) == 10 {
		short := new(models.Short)

		short.Short = request.URL.Path[1:]

		if short.Short == "" {
			response.Failed(writer, http.StatusBadRequest, errors.New("Bad payload"))
			return
		}

		if err := short.Get(); err != nil {
			response.Failed(writer, http.StatusNotFound, err)
			return
		}
		response.Success(writer, http.StatusOK, short)
		return
	}
	response.Failed(writer, http.StatusBadRequest, errors.New("Bad url path"))
	return
}
