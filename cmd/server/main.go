package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"maxim.tbank/xml2pg/conc"
	"maxim.tbank/xml2pg/db"
	"maxim.tbank/xml2pg/handlers"
)

type VarcharFunc func(ctx context.Context, msg string) error

func (f VarcharFunc) InsertMessage(ctx context.Context, msg string) error {
	return f(ctx, msg)
}
func (f VarcharFunc) InsertMessageJson(ctx context.Context, msg string) error {
	return f(ctx, msg)
}

func main() {
	var (
		serverHost  = GetEnv("SERVER_HOST", "0.0.0.0")
		serverPort  = GetEnv("SERVER_PORT", "8080")
		logLevel    = GetEnv("LOG_LEVEL", "INFO")
		postgresUrl = GetEnv("POSTGRES_URL", "postgres://postgres:postgres_pass@localhost:5432/maxim?pool_max_conns=10")
	)
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: false}).With().Timestamp().Logger()

	ConfigureLogLevelOrDie(logLevel)
	log.Info().Msg("start xml2pg")

	log.Info().Msgf("NumCPU: %d", runtime.NumCPU())
	log.Info().Msgf("LOG_LEVEL: %s", zerolog.GlobalLevel())
	log.Info().Msgf("SERVER_HOST: %s", serverHost)
	log.Info().Msgf("SERVER_PORT: %s", serverPort)
	uri := ParsePostgresUrlOrDie(postgresUrl)
	log.Info().Msgf("POSTGRES_URL: %s", uri.Redacted())

	globalCtx, globalContextCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer globalContextCancel()

	ctx, c := context.WithTimeout(globalCtx, 10*time.Second)
	defer c()
	pool, err := db.NewPostgresPool(ctx, postgresUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create postgres pool")
		return
	}
	storage := db.NewPostgresStorage(pool)

	errorFunc := func(e error) {
		log.Error().Err(e).Msg("Batch error")
	}
	batchVarchar := conc.NewBatchf[string](errorFunc)

	go batchVarchar.CollectAndExec(globalCtx, 32, 50*time.Millisecond, storage.InsertMessages)

	batchJsonb := conc.NewBatchf[string](errorFunc)

	go batchJsonb.CollectAndExec(globalCtx, 32, 50*time.Millisecond, storage.InsertMessagesJson)

	mux := http.NewServeMux()
	mux.HandleFunc("/health/live", handlers.Healthcheck)

	mux.HandleFunc("POST /message/varchar", handlers.MessageVarcharInsert(storage))
	mux.HandleFunc("POST /batch/varchar", handlers.MessageVarcharInsert(VarcharFunc(batchVarchar.Execute)))

	mux.HandleFunc("POST /message/jsonb", handlers.MessageJsonInsert(storage))
	mux.HandleFunc("POST /batch/jsonb", handlers.MessageJsonInsert(VarcharFunc(batchJsonb.Execute)))

	srvShutdownFunc := startHTTPServerOrDie(net.JoinHostPort(serverHost, serverPort), mux)

	<-globalCtx.Done()

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	errSrv := srvShutdownFunc(ctx)
	if errSrv != nil {
		log.Error().Err(errSrv).Msg("error shutdown http server")
	}

}

func GetEnv(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func ConfigureLogLevelOrDie(level string) {
	logLvl, err := zerolog.ParseLevel(level)
	if err != nil {
		log.Fatal().Err(err).Msg("Error parsing log level")
		return
	}
	zerolog.SetGlobalLevel(logLvl)
}

func ParsePostgresUrlOrDie(postgresUrl string) *url.URL {
	uri, err := url.Parse(postgresUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("error parse postgres url")
		return nil
	}
	return uri
}

type ShutdownFunc func(ctx context.Context) error

func startHTTPServerOrDie(addr string, handler http.Handler) ShutdownFunc {
	s := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	go func() {
		err := s.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Info().Msgf("server closed: %s", addr)
		} else if err != nil {
			log.Fatal().Err(err).Msg("Unable to create http server")
		}
	}()

	return func(ctx context.Context) error {
		return s.Shutdown(ctx)
	}
}
