package webserver

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/7onn/croupier/internal/croupier"
	"github.com/rs/zerolog/log"
)

type Server interface {
	croupier.Croupier
}

// Broker manages the internal state of the Croupier API.
type Broker struct {
	croupier.Croupier

	cfg    Config      // the api service's configuration
	router *mux.Router // the api service's route collection
}

// New initializes a new Croupier API.
func New(cfg Config) (*Broker, error) {
	r := &Broker{}

	err := validateConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}
	r.cfg = cfg

	r.Croupier = croupier.New()

	// Do other setup work here...

	return r, nil
}

// Start the Croupier service
func (bkr *Broker) Start(binder func(s Server, r *mux.Router)) {
	bkr.router = mux.NewRouter().StrictSlash(true)
	binder(bkr, bkr.router)

	// Do other startup work here...
	l, err := net.Listen("tcp", ":"+strconv.Itoa(bkr.cfg.Port))
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to bind to TCP port %d for listening.", bkr.cfg.Port)
		os.Exit(13)
	} else {
		log.Info().Msgf("Starting webserver on TCP port %04d", bkr.cfg.Port)
	}

	if err := http.Serve(l, bkr.router); errors.Is(err, http.ErrServerClosed) {
		log.Warn().Err(err).Msg("Web server has shut down")
	} else {
		log.Fatal().Err(err).Msg("Web server has shut down unexpectedly")
	}
}
