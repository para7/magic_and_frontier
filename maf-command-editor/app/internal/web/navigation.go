package web

import (
	"net/http"
	neturl "net/url"
	"strings"

	"tools2/app/internal/webui"
)

func currentListURL(r *http.Request, fallback string) string {
	return sanitizeReturnTo(r.URL.RequestURI(), fallback)
}

func queryReturnTo(r *http.Request, fallback string) string {
	return sanitizeReturnTo(r.URL.Query().Get("returnTo"), fallback)
}

func submittedReturnTo(r *http.Request, fallback string) string {
	return sanitizeReturnTo(r.Form.Get("returnTo"), fallback)
}

func applyPageMeta(r *http.Request, meta webui.PageMeta) webui.PageMeta {
	if strings.TrimSpace(meta.CurrentURL) == "" {
		meta.CurrentURL = currentListURL(r, meta.CurrentPath)
	}
	return meta
}

func sanitizeReturnTo(value string, fallback string) string {
	fallback = normalizeScreenPath(fallback)
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return fallback
	}
	parsed, err := neturl.Parse(trimmed)
	if err != nil {
		return fallback
	}
	if parsed.Scheme != "" || parsed.Host != "" || parsed.User != nil {
		return fallback
	}
	if !strings.HasPrefix(parsed.Path, "/") {
		return fallback
	}
	if parsed.Path != fallback {
		return fallback
	}
	parsed.Scheme = ""
	parsed.Host = ""
	parsed.User = nil
	parsed.Fragment = ""
	result := parsed.EscapedPath()
	if result == "" {
		result = fallback
	}
	if parsed.RawQuery != "" {
		result += "?" + parsed.RawQuery
	}
	return result
}
