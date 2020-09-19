package utils

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

// testCert self generated public cert for testing
const testCert = `
-----BEGIN CERTIFICATE-----
MIIC5TCCAc2gAwIBAgIJAN47QcKwvWpPMA0GCSqGSIb3DQEBCwUAMBQxEjAQBgNV
BAMMCWxvY2FsaG9zdDAeFw0yMDA5MTkxODExMzVaFw0yMDEwMTkxODExMzVaMBQx
EjAQBgNVBAMMCWxvY2FsaG9zdDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoC
ggEBAJS/32WujYi0ZV8p+7uF2aXWbc5Kxb/qnr1AqEFgVpgS2TaVyxsiajA6vSbR
6a5/OCeKuneUrd+M92jWxyfR3/vDy0y7G8mSeUhM+Ln4aXIV63F0pUApo14jCQHt
/SXJbhsbuFUz0sXDqviOTGtjJg4ymaTp9gAoF4UCTy069haH9aR4p9rFwMUr0unF
Pzzehr5p3bl934JIOcimUwzpvYLavgxzN0O3y+yULa4v6lSZzsdz0wCkaXFud1Ky
JfFDonKqdvP3ui1pPfPrLT2E4tD1KeLKdi8a0MlAGPLgdjgonboNtg6BQApF6uTf
vTB1hTv+bcpWtAlZJXI15VtN01ECAwEAAaM6MDgwFAYDVR0RBA0wC4IJbG9jYWxo
b3N0MAsGA1UdDwQEAwIHgDATBgNVHSUEDDAKBggrBgEFBQcDATANBgkqhkiG9w0B
AQsFAAOCAQEAOp+71j4HWI+zBbDmINlJ3dpI+ZVWAPvpjTM54IRW1yXlAbYfHQPb
p3VCXVc3fuVwxlWUfyMpeokSyibkNfZHENHSd+TnTAzb5B9UkQ2BlGks19kn3IIh
PdJVJ2lRX1qgjxBrNCzJg6thyBH0mOCv1y+fizbl+F7N/6Eosw/X0HfQx+spaSZY
W0MtXxfaELY0IkI3W2tRmr+etxREW9u2qZPhn9YyZJLP09h6W0m4ksje9E4oDTSk
vQgEfl+93ZLZQs6FcJswGmmtjrRhmEvagUxr8ftpwnir7z/I3MmonmPbymjC5DSZ
3YUp3pVRkTL5S7P7l6oVaAGP4Ld1AYC12Q==
-----END CERTIFICATE-----
`

// Client http client to make requests
type Client struct {
	client *http.Client
}

// NewClient construct http client
func NewClient(timeout time.Duration) *Client {
	return &Client{client: &http.Client{
		Transport: transport(),
		Timeout:   timeout,
	}}
}

// tlsConfig read public cert and return tls.Config
//nolint:gosec
func tlsConfig() *tls.Config {
	rootCAs := x509.NewCertPool()
	rootCAs.AppendCertsFromPEM([]byte(testCert))

	return &tls.Config{
		RootCAs:            rootCAs,
		InsecureSkipVerify: true,
	}
}

// transport configuration for tls
func transport() *http.Transport {
	return &http.Transport{
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		DisableKeepAlives:     true,
		TLSClientConfig:       tlsConfig(),
		DisableCompression:    true,
	}
}

// NewTestTLSServer start test server over TLS
func NewTestTLSServer(h http.Handler) *httptest.Server {
	srv := httptest.NewUnstartedServer(h)
	srv.TLS = tlsConfig()
	srv.StartTLS()
	return srv
}

// Get simple get request
func (c *Client) Get(ctx context.Context, u *url.URL, successResponse, errorResponse interface{}) (int, error) {
	return c.GetWithHeaders(ctx, u, successResponse, errorResponse, nil)
}

// GetWithHeaders simple get request with headers
//nolint:interfacer
func (c *Client) GetWithHeaders(ctx context.Context, u *url.URL, successResponse, errorResponse interface{}, headers map[string][]string) (int, error) {
	if u == nil {
		return -1, errors.New("url shouldn't be nil")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return -1, errors.WithStack(err)
	}

	// set headers if any
	if headers != nil {
		req.Header = headers
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return -1, errors.WithStack(err)
	}
	defer resp.Body.Close()

	jd := json.NewDecoder(resp.Body)

	if errorResponse != nil && resp.StatusCode >= http.StatusMultipleChoices {
		return resp.StatusCode, errors.WithStack(jd.Decode(errorResponse))
	}

	if successResponse != nil {
		return resp.StatusCode, errors.WithStack(jd.Decode(successResponse))
	}

	return resp.StatusCode, nil
}
