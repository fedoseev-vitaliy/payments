package gpay

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// MockGPay is mock for like e2e test
type MockGPay struct{}

func (mg *MockGPay) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pid := r.URL.Query().Get("productID")
	if pid == "" {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(&googlePayError{
			Error: "orderID query params is missing",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if pid == "badApple" {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(&googlePayError{
			Error: "bad apple product",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(&googlePayResponse{
		PayButtonURL: fmt.Sprintf("http://google.pay.com/payfor?product=%s", pid),
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
