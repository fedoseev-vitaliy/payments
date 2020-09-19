package apay

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// MockAPay is mock for like e2e test
type MockAPay struct{}

func (ma *MockAPay) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pid := r.URL.Query().Get("productID")
	if pid == "" {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(&applePayError{
			Error: "orderID query params is missing",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if pid == "badGoogle" {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(&applePayError{
			Error: "bad google product",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(&applePayResponse{
		PayButtonURL: fmt.Sprintf("http://apple.pay.com/payfor?product=%s", pid),
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
