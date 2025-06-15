package main

import (
	"context"
	"encoding/json"
	"math/rand/v2"
	"net/http"
)

type JSONAPIServer struct {
	listenAddr	string
	svc 				PriceFetcher
}

type PriceResponse struct {
	Ticker 	string		`json:"ticker"`
	Price		float64		`json:"price"`
}

type APIFunc func(context.Context, http.ResponseWriter, *http.Request) error

func newJSONAPIServer(listenAddr string, svc PriceFetcher) *JSONAPIServer {
	return &JSONAPIServer{
		listenAddr: listenAddr,
		svc: svc,
	}
}

func (s *JSONAPIServer) Run() {
	http.HandleFunc("/", makeHTTPHandlerFunc(s.handleFetchPrice));
	http.ListenAndServe(s.listenAddr, nil);
}

func makeHTTPHandlerFunc(apiFn APIFunc) http.HandlerFunc {
	ctx := context.Background();
	ctx = context.WithValue(ctx, "requestId", rand.IntN(10000000));

	return func(w http.ResponseWriter, r *http.Request) {
		if err := apiFn(ctx, w, r); err != nil {
			w.WriteHeader(http.StatusBadRequest);
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()});
		}
	}
}

func (s *JSONAPIServer) handleFetchPrice(
	ctx context.Context, 
	w http.ResponseWriter,
	r *http.Request,
) error {
	ticker := r.URL.Query().Get("ticker");

	price, err := s.svc.FetchPrice(ctx, ticker);
	if err != nil {
		return err;
	}

	priceResp := PriceResponse{
		Price: price,
		Ticker: ticker,
	}

	return writeJSON(w, http.StatusOK, &priceResp);
}

func writeJSON(w http.ResponseWriter, s int, v any) error {
	w.WriteHeader(s);
	return json.NewEncoder(w).Encode(v);
}