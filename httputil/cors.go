package httputil

import "net/http"

type (
	CorsConfig struct {
		Origins []string
		Methods []string
		Headers []string
	}
)

// Access-Control-Allow-Origin defaults to "*" if no origins are configured.
func CorsHandlerFunc(config CorsConfig, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(config.Origins) < 1 {
			w.Header().Add("access-control-allow-origin", "*")
		} else {
			for _, origin := range config.Origins {
				if origin == r.Header.Get("origin") {
					w.Header().Add("access-control-allow-origin", origin)
					break
				}
			}
		}

		for _, method := range config.Methods {
			w.Header().Add("access-control-allow-methods", method)
		}

		for _, header := range config.Headers {
			w.Header().Add("access-control-allow-headers", header)
		}

		next(w, r)
	}
}

func CorsHandler(config CorsConfig, next http.Handler) http.Handler {
	return CorsHandlerFunc(config, next.ServeHTTP)
}
