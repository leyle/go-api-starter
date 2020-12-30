package logmiddleware

import (
	"fmt"
	"github.com/rs/zerolog"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	logger := zerolog.Ctx(r.Context())

	logger.Info().Msg("Start processing...")

	_, err := w.Write([]byte("OK"))
	if err != nil {
		fmt.Println(err)
	}
	logger.Info().Str("id", "userapp01").Str("name", "admin").Msg("")
	body, _ := ioutil.ReadAll(r.Body)
	logger.Debug().RawJSON("body", body).Msg("")

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
	// logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	// logger := GetLogger(LogTargetConsole)
	logger := GetLogger(LogTargetStdout)

	mux := http.NewServeMux()
	mux.Handle("/", ZeroLogMiddleware(logger, http.HandlerFunc(testHandler)))

	err := http.ListenAndServe(":8080", mux)
	logger.Fatal().Err(err)
}
