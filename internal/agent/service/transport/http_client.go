package transport

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/dontagr/metric/internal/agent/config"
)

type HTTPManager struct {
	client *http.Client
	log    *zap.SugaredLogger
	url    string
	ip     string
}

func NewHTTPManager(cfg *config.Config, log *zap.SugaredLogger) (*HTTPManager, error) {
	httpManager := HTTPManager{log: log, client: &http.Client{}}
	if cfg.RateLimit == 0 {
		httpManager.url = fmt.Sprintf("http://%s/updates/", cfg.HTTPBindAddress)
	} else {
		httpManager.url = fmt.Sprintf("http://%s/update/", cfg.HTTPBindAddress)
	}

	ip, err := getIP()
	if err != nil {
		return nil, fmt.Errorf("failed get ip: %v", err)
	}
	httpManager.ip = ip

	return &httpManager, nil
}

func (h *HTTPManager) NewRequest(income any, HashSHA256 []string, w int) error {
	compressedBody, ok := income.(*bytes.Buffer)
	if !ok {
		return fmt.Errorf("failed to convert data to *bytes.Buffer")
	}

	req, err := http.NewRequest("POST", h.url, compressedBody)
	if err != nil {
		return fmt.Errorf("creating request: %v", err)
	}

	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Real-IP", h.ip)
	for _, hashRow := range HashSHA256 {
		req.Header.Add("HashSHA256", hashRow)
	}

	var resp *http.Response
	var netErr *net.OpError
	var errSend error
	for i := 0; i < 3; i++ {
		resp, errSend = h.client.Do(req) // nolint
		if errSend == nil {
			statusCode := resp.StatusCode
			if statusCode >= 200 && statusCode < 300 {
				h.log.Infof("worker %d request sent successfully with status code: %d", w, statusCode)
				return nil
			} else {
				h.log.Warnf("worker %d received non-2xx status code: %d", w, statusCode)
				errSend = fmt.Errorf("received non-2xx status code: %d", statusCode)
			}

			resp.Body.Close()
			return errSend
		}
		if errors.As(errSend, &netErr) {
			h.log.Warnf("worker %d connection error we try №%d", w, i+1)
			if i < 2 {
				time.Sleep(5 * time.Second)
			}
		} else {
			return fmt.Errorf("sending data: %v", errSend)
		}
	}

	return fmt.Errorf("failed to send request after retrying: %v", errSend)
}

func getIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("failed in InterfaceAddrs: %v", err)
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "127.0.0.1", nil
}
