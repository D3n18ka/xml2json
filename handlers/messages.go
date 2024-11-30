package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"

	"maxim.tbank/xml2pg/service"
)

type OneMessageInserter interface {
	InsertMessage(ctx context.Context, msg string) error
}

func MessageVarcharInsert(db OneMessageInserter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.Err(err).Msg("unable to read body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		msg, err := service.ParseXmlStruct(data)
		if err != nil {
			log.Err(err).Msg("ParseXmlStruct")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		marshal, err := json.Marshal(msg)
		if err != nil {
			log.Err(err).Msg("json.Marshal")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = db.InsertMessage(context.Background(), string(marshal))
		if err != nil {
			log.Err(err).Msg("InsertMessage")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
