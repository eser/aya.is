package profiling

import (
	"net/http/pprof"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
)

func RegisterHTTPRoutes(routes *httpfx.Router, config *httpfx.Config) {
	if !config.ProfilingEnabled {
		return
	}

	mux := routes.GetMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}
