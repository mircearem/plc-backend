package Helmet

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Config struct {
	// Filter defines a function to skip middleware.
	// Optional. Default: nil
	Filter func(w http.ResponseWriter, r *http.Request) bool
	// XSSProtection
	// Optional. Default value "1; mode=block".
	XSSProtection string
	// ContentTypeNosniff
	// Optional. Default value "nosniff".
	ContentTypeNosniff string
	// XFrameOptions
	// Optional. Default value "SAMEORIGIN".
	// Possible values: "SAMEORIGIN", "DENY", "ALLOW-FROM uri"
	XFrameOptions string
	// HSTSMaxAge
	// Optional. Default value 0.
	HSTSMaxAge int
	// HSTSExcludeSubdomains
	// Optional. Default value false.
	HSTSExcludeSubdomains bool
	// ContentSecurityPolicy
	// Optional. Default value "".
	ContentSecurityPolicy string
	// CSPReportOnly
	// Optional. Default value false.
	CSPReportOnly bool
	// HSTSPreloadEnabled
	// Optional.  Default value false.
	HSTSPreloadEnabled bool
	// ReferrerPolicy
	// Optional. Default value "".
	ReferrerPolicy string

	// Permissions-Policy
	// Optional. Default value "".
	PermissionPolicy string
}

func New(config ...Config) mux.MiddlewareFunc {
	var cfg Config
	if len(config) > 0 {
		cfg = config[0]
	}
	// Set config default values
	if cfg.XSSProtection == "" {
		cfg.XSSProtection = "1; mode=block"
	}
	if cfg.ContentTypeNosniff == "" {
		cfg.ContentTypeNosniff = "nosniff"
	}
	if cfg.XFrameOptions == "" {
		cfg.XFrameOptions = "SAMEORIGIN"
	}

	// Return middleware handler
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cfg.XSSProtection != "" {
				w.Header().Add("X-XSS-Protection", cfg.XSSProtection)
			}
			if cfg.ContentTypeNosniff != "" {
				w.Header().Add("X-Content-Type-Options", cfg.ContentTypeNosniff)
			}
			if cfg.XFrameOptions != "" {
				w.Header().Add("X-Frame-Options", cfg.XFrameOptions)
			}

			if (r.Header.Get("X-Forwarder-Proto") == "https") && cfg.HSTSMaxAge != 0 {
				subdomains := ""

				if !cfg.HSTSExcludeSubdomains {
					subdomains = "; includeSubdomains"
				}

				if cfg.HSTSPreloadEnabled {
					subdomains = fmt.Sprintf("%s; preload", subdomains)
				}

				w.Header().Add("Strict-Transport-Security", fmt.Sprintf("max-age=%d%s", cfg.HSTSMaxAge, subdomains))
			}

			if cfg.ContentSecurityPolicy != "" {
				if cfg.CSPReportOnly {
					w.Header().Add("Content-Security-Policy-Report-Only", cfg.ContentSecurityPolicy)
				} else {
					w.Header().Add("Content-Security-Policy", cfg.ContentSecurityPolicy)
				}
			}

			if cfg.ReferrerPolicy != "" {
				w.Header().Add("Referrer-Policy", cfg.ReferrerPolicy)
			}
			if cfg.PermissionPolicy != "" {
				w.Header().Add("Permissions-Policy", cfg.PermissionPolicy)
			}

			h.ServeHTTP(w, r)
		})
	}
}
