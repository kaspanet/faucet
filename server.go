package main

import (
	"context"
	"net/http"
	"time"

	"github.com/kaspanet/faucet/config"
	"github.com/kaspanet/kaspad/util"
	"github.com/kaspanet/kasparov/httpserverutils"
	"github.com/pkg/errors"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const gracefulShutdownTimeout = 30 * time.Second

// startHTTPServer starts the HTTP REST server and returns a
// function to gracefully shutdown it.
func startHTTPServer(listenAddr string) func() {
	router := mux.NewRouter()
	router.Use(httpserverutils.AddRequestMetadataMiddleware)
	router.Use(httpserverutils.RecoveryMiddleware)
	router.Use(httpserverutils.LoggingMiddleware)
	router.Use(httpserverutils.SetJSONMiddleware)
	router.HandleFunc(
		"/request_money?address={address}",
		httpserverutils.MakeHandler(requestMoneyHandler)).
		Methods("GET")
	httpServer := &http.Server{
		Addr:    listenAddr,
		Handler: handlers.CORS()(router),
	}
	spawn("startHTTPServer-httpServer.ListenAndServe", func() {
		log.Errorf("%s", httpServer.ListenAndServe())
	})

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
		defer cancel()
		err := httpServer.Shutdown(ctx)
		if err != nil {
			log.Errorf("Error shutting down HTTP server: %s", err)
		}
	}
}

func requestMoneyHandler(_ *httpserverutils.ServerContext, request *http.Request,
	routeParams map[string]string, _ map[string]string, _ []byte) (interface{}, error) {

	hErr := validateIPUsage(request)
	if hErr != nil {
		return nil, hErr
	}
	addressString, ok := routeParams["address"]
	if !ok {
		return nil, httpserverutils.NewHandlerErrorWithCustomClientMessage(http.StatusUnprocessableEntity,
			errors.Errorf("address not found"),
			"The address parameter is either missing or empty")
	}
	address, err := util.DecodeAddress(addressString, config.ActiveNetParams().Prefix)
	if err != nil {
		return nil, httpserverutils.NewHandlerErrorWithCustomClientMessage(http.StatusUnprocessableEntity,
			errors.Wrap(err, "Error decoding address"),
			"Error decoding address")
	}
	transactionID, err := sendToAddress(address)
	if err != nil {
		return nil, err
	}
	hErr = updateIPUsage(request)
	if hErr != nil {
		return nil, hErr
	}
	return transactionID, nil
}
