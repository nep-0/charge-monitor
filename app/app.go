package app

import (
	"charge-monitor/cache"
	"charge-monitor/query"
	"log/slog"
	"net/http"
	"time"
)

type App struct {
	outlets         []string
	pollingInterval time.Duration
	httpAddress     string
	cache           cache.Cache
}

func NewApp(outlets []string, pollingInterval int64, httpAddress string) *App {
	return &App{
		outlets:         outlets,
		pollingInterval: time.Duration(pollingInterval) * time.Millisecond,
		httpAddress:     httpAddress,
		cache:           cache.NewLocalCache(),
	}
}

func (a *App) ServeHTTP() {
	go a.poll()
	http.HandleFunc("/outlets", a.corsMiddleware(a.getOutlets))
	slog.Info("Starting HTTP server", "address", a.httpAddress)
	http.ListenAndServe(a.httpAddress, nil)
}

func (a *App) getOutlets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(a.cache.JSON())
}

func (a *App) corsMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler(w, r)
	}
}

func (a *App) poll() {
	for {
		errorCount := 0
		for _, outletId := range a.outlets {
			power, usedMinutes, err := query.QueryChargeStatus(outletId)
			if err != nil {
				slog.Error("Failed to query charge status", "outletId", outletId, "error", err)
				errorCount++
				continue
			}
			a.cache.Set(outletId, cache.OutletInfo{Power: power, UsedMinutes: usedMinutes})
			time.Sleep(a.pollingInterval)
		}
		slog.Info("Completed a full polling cycle", "errors", errorCount)
	}
}
