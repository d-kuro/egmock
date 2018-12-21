package serve

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/d-kuro/egmock/log"
)

type Mock struct {
	status  int
	resBody string
}

type RequestLog struct {
	Protocol    string `json:"protocol"`
	ContentType string `json:"content_type"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	Query       string `json:"query"`
	Body        string `json:"body"`
}

func NewMock(status int, resBody string) *Mock {
	return &Mock{
		status:  status,
		resBody: resBody,
	}
}

func (m *Mock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// request logging
	bufBody := new(bytes.Buffer)
	_, err := io.Copy(bufBody, r.Body)
	if err != nil {
		log.Error("get request body error", zap.Error(err))
		w.WriteHeader(500)
	}

	reqLog := RequestLog{
		Protocol:    r.Proto,
		ContentType: r.Header.Get("Content-Type"),
		Method:      r.Method,
		Path:        r.URL.Path,
		Query:       r.URL.RawQuery,
		Body:        bufBody.String(),
	}

	jsonBytes, err := json.Marshal(reqLog)
	if err != nil {
		log.Error("json marshal error", zap.Error(err))
		w.WriteHeader(500)
	}
	log.Info(string(jsonBytes))

	w.WriteHeader(m.status)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, m.resBody)
}
