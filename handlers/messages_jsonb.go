package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"maxim.tbank/xml2pg/service"
)

type OneMessageJsonBInserter interface {
	InsertMessageJson(ctx context.Context, msg string) error
}

func MessageJsonInsert(db OneMessageJsonBInserter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		msg, err := service.ParseXmlStruct(data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		marshal, err := json.Marshal(msg)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = db.InsertMessageJson(r.Context(), string(marshal))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
