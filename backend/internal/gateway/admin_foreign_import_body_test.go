package gateway

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// zeroReader returns infinite zero bytes without allocating a large buffer.
type zeroReader struct{}

func (zeroReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}

func TestAdminImportForeignConfigBodyTooLarge(t *testing.T) {
	s, _ := newBulkTestServer(t)

	// Stream slightly more than the cap without allocating 64 MiB.
	body := io.LimitReader(zeroReader{}, foreignImportMaxBytes+1)
	req := httptest.NewRequest(http.MethodPost, "/admin/import/foreign", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.adminImportForeignConfig(rec, req)

	require.Equal(t, http.StatusRequestEntityTooLarge, rec.Code)
	msg := rec.Body.String()
	require.True(t,
		strings.Contains(msg, "too large") || strings.Contains(msg, "64") || strings.Contains(msg, "MiB"),
		"expected oversized body message, got %q", msg)
}
