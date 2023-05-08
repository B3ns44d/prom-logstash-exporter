package restclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

type HTTPHandler struct {
	Endpoint string
}

func (h *HTTPHandler) Get() (*http.Response, error) {
	response, err := http.Get(h.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to GET %s: %v", h.Endpoint, err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s returned status code %d", h.Endpoint, response.StatusCode)
	}

	return response, nil
}

type HTTPHandlerInterface interface {
	Get() (*http.Response, error)
}

func GetMetrics(h HTTPHandlerInterface, target interface{}) error {
	response, err := h.Get()
	if err != nil {
		return fmt.Errorf("failed to retrieve metrics data: %v", err)
	}
	defer closeResponseBody(response.Body)

	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to decode metrics data: %v", err)
	}

	return nil
}

func closeResponseBody(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		logrus.Errorf("cannot close response body: %v", err)
	}
}
