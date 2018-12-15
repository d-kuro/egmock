package serve

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/d-kuro/egmock/logger"
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
		logger.ELog.Println("get request body error:", err)
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
		logger.ELog.Println("json marshal error:", err)
		w.WriteHeader(500)
	}
	logger.ILog.Println(string(jsonBytes))

	w.WriteHeader(m.status)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write([]byte(m.resBody))
	if err != nil {
		logger.ELog.Println("write response body error:", err)
		w.WriteHeader(500)
	}
}
