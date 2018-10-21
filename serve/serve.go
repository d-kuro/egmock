package serve

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/d-kuro/egmock/logger"
)

type mock struct {
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

func NewMock(status int, resBody string) *mock {
	return &mock{
		status:  status,
		resBody: resBody,
	}
}

func (m *mock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// request logging
	bufBody := new(bytes.Buffer)
	io.Copy(bufBody, r.Body)

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
	}
	logger.ILog.Println(string(jsonBytes))

	w.WriteHeader(m.status)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(m.resBody))
	return
}
