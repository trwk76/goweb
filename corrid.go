package goweb

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func CorrID(r *http.Request) string {
	return r.Header.Get(corrIDHeader)
}

func newCorrID() string {
	return strings.ToUpper(strings.ReplaceAll(uuid.NewString(), "-", ""))
}

func setCorrID(r *http.Request, corrID string) {
	r.Header.Set(corrIDHeader, corrID)
}

const corrIDHeader string = "X-CorrID"

func init() {
	uuid.EnableRandPool()
}
