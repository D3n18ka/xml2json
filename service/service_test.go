package service_test

import (
	"context"
	_ "embed"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/rs/zerolog/log"

	"maxim.tbank/xml2pg/db"
)

//go:embed data.json
var dataJson string

func TestDb(t *testing.T) {

	postgresUrl := "postgres://postgres:postgres_pass@localhost:5432/maxim"

	globalCtx, globalContextCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer globalContextCancel()

	ctx, c := context.WithTimeout(globalCtx, 10*time.Second)
	defer c()
	pool, err := db.NewPostgresPool(ctx, postgresUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create postgres pool")
		return
	}
	storage := db.NewPostgresStorage(pool)

	//err = storage.InsertMessage(globalCtx, "test message")
	//if err != nil {
	//	t.Error("error inserting message", err)
	//}

	msgs := make([]string, 0)
	for i := range 1024 {
		msgs = append(msgs, "test message "+strconv.Itoa(i))
	}

	err = storage.InsertMessages(globalCtx, msgs)
	if err != nil {
		t.Error("error inserting messages", err)
	}

	//err = storage.InsertMessageJson(globalCtx, dataJson)
	//if err != nil {
	//	t.Error("error inserting message", err)
	//}

}
