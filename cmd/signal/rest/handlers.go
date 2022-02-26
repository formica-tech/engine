package rest

import (
	"encoding/json"
	"fmt"
	"github.com/hamzali/formica-engine/domain"
	"github.com/hamzali/formica-engine/usecases"
	"io"
	"io/ioutil"
	"net/http"
)

type Handler struct {
	signalUseCases usecases.SignalUseCases
}

type Response struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func (h Handler) respondMsg(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	respData := Response{Message: message}
	resp, _ := json.Marshal(&respData)
	_, _ = w.Write(resp)
}

func (h Handler) respondData(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	respData := Response{Data: data}
	resp, err := json.Marshal(&respData)
	if err != nil {
		h.respondMsg(w, http.StatusInternalServerError, "could not parse response body")
		return
	}
	_, _ = w.Write(resp)
}

func (h Handler) respondUnexpected(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	h.respondMsg(w, http.StatusInternalServerError, err.Error())
}

func (h Handler) parseBody(w http.ResponseWriter, body io.ReadCloser, data interface{}) error {
	bodyStr, err := ioutil.ReadAll(body)
	if err != nil {
		h.respondMsg(w, http.StatusBadRequest, "could not read request body")
		return err
	}
	err = json.Unmarshal(bodyStr, data)
	if err != nil {
		h.respondMsg(w, http.StatusBadRequest, "could not parse request body")
		return err
	}
	return nil
}

func (h Handler) BatchInsertSdkSignal(w http.ResponseWriter, r *http.Request) {

	signals := make([]domain.EntitySignal, 0)
	err := h.parseBody(w, r.Body, &signals)
	if err != nil {
		return
	}

	if len(signals) == 0 {
		h.respondMsg(w, http.StatusBadRequest, "signal batch can't be empty")
		return
	}

	err = h.signalUseCases.BatchSave(r.Context(), signals)
	if err != nil {
		h.respondUnexpected(w, err)
		return
	}

	h.respondMsg(w, http.StatusOK, fmt.Sprintf("%d signals are saved", len(signals)))
}
