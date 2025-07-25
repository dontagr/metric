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
	url    string
	client *http.Client
	log    *zap.SugaredLogger
}

func NewHTTPManager(cfg *config.Config, log *zap.SugaredLogger) *HTTPManager {
	httpManager := HTTPManager{log: log, client: &http.Client{}}
	if cfg.RateLimit == 0 {
		httpManager.url = fmt.Sprintf("http://%s/updates/", cfg.HTTPBindAddress)
	} else {
		httpManager.url = fmt.Sprintf("http://%s/update/", cfg.HTTPBindAddress)
	}

	return &httpManager
}

func (h *HTTPManager) NewRequest(compressedBody *bytes.Buffer, HashSHA256 []string, w int) error {
	req, err := http.NewRequest("POST", h.url, compressedBody)
	if err != nil {
		return fmt.Errorf("creating request: %v", err)
	}

	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	for _, hashRow := range HashSHA256 {
		req.Header.Add("HashSHA256", hashRow)
	}

	var resp *http.Response
	var netErr *net.OpError
	var errSend error
	bodyClose := false
	for i := 0; i < 3; i++ {
		resp, errSend = h.client.Do(req)
		if errSend == nil {
			// just for linter
			err = resp.Body.Close()
			if err != nil {
				return fmt.Errorf("closing response body: %v", err)
			}
			h.log.Infof("worker %d request success full", w)
			bodyClose = true
			return nil
		}
		if errors.As(errSend, &netErr) {
			h.log.Warnf("worker %d connection error we try №%d", w, i+1)
			time.Sleep(5 * time.Second)
		} else {
			return fmt.Errorf("sending data: %v", errSend)
		}
	}

	if errSend != nil {
		return fmt.Errorf("sending data: %v", errSend)
	}

	if !bodyClose {
		err = resp.Body.Close()
		if err != nil {
			return fmt.Errorf("closing response body: %v", err)
		}
		h.log.Infof("worker %d request success full", w)
	}

	return nil
}
