package tripper

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

type RetryRoundTripper struct {
	tripper http.RoundTripper
	retryOn []int
	times   int
	waitFor time.Duration
}

type RetryConfig struct {
	RetryOn   []int
	Times     int
	WaitForMs int
}

func NewRetryRoundTripper(tripper http.RoundTripper, config RetryConfig) *RetryRoundTripper {
	return &RetryRoundTripper{
		tripper: tripper,
		retryOn: config.RetryOn,
		times:   config.Times,
		waitFor: time.Duration(config.WaitForMs) * time.Microsecond,
	}
}

func (r *RetryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}

		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	resp, err := r.tripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	for i := 0; i < r.times && isErrorRetry(resp, r.retryOn); i++ {
		time.Sleep(r.waitFor)

		drainResponseBody(resp)

		if req.Body != nil {
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		resp, err = r.tripper.RoundTrip(req)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func drainResponseBody(resp *http.Response) {
	if resp.Body != nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func isErrorRetry(response *http.Response, retryOn []int) bool {
	for _, code := range retryOn {
		if response.StatusCode == code {
			return true
		}
	}
	return false
}
