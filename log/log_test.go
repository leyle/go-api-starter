package log

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	logger := log.Ctx(r.Context())

	logger.Info().Msg("Start processing...")

	_, err := w.Write([]byte("OK"))
	if err != nil {
		fmt.Println(err)
	}
	logger.Info().Str("id", "userapp01").Str("name", "admin").Msg("")
	body, _ := ioutil.ReadAll(r.Body)
	log.Debug().RawJSON("body", body).Msg("")

	logger.Info().Dur("elapsed", time.Since(start)).Msg("Done")
}

func TestLogMiddleware(t *testing.T) {
	// zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	/*
		logger := zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}).With().Timestamp().Logger()
	*/
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	mux := http.NewServeMux()
	mux.Handle("/",
		LogMiddleware(logger)(
			ReqIdMiddleware(
				http.HandlerFunc(testHandler),
			),
		),
	)

	err := http.ListenAndServe(":8080", mux)
	log.Fatal().Err(err)
}
