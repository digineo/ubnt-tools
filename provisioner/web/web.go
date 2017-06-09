package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/digineo/ubnt-tools/provisioner"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	headerContentLength   = "Content-Length"
	headerContentType     = "Content-Type"
	headerIfModifiedSince = "If-Modified-Since"
	headerLastModified    = "Last-Modified"
	httpDateFormat        = time.RFC1123
	cacheDuration         = 20 * time.Second

	contentTypeHTML = "text/html; charset=UTF-8"
	contentTypeJSON = "application/json; charset=UTF-8"
)

type goWeb struct {
	config *provisioner.Configuration
	router *mux.Router
	server *http.Server
}

// StartWeb initializes the web UI
func StartWeb(c *provisioner.Configuration) {
	web := &goWeb{config: c}
	web.buildRoutes()

	cors := handlers.CORS(handlers.AllowedOrigins([]string{"*"}))
	handler := handlers.LoggingHandler(os.Stdout, cors(web.router))

	web.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", c.Web.Host, c.Web.Port),
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Printf("[webui] Starting Web UI on http://%s", web.server.Addr)
	web.server.ListenAndServe()
}

func (g *goWeb) buildRoutes() {
	g.router = mux.NewRouter()

	g.router.HandleFunc("/", g.getRoot).Methods("GET").Name("root")
	g.router.PathPrefix("/assets/").HandlerFunc(g.getStaticAsset).Name("asset")
	g.router.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		routes := make(map[string]string)
		for _, name := range []string{"api_directory", "device", "devices", "upgrade_device", "provision_device", "reboot_device"} {
			route := g.router.Get(name)
			if tpl, err := route.GetPathTemplate(); err == nil {
				routes[name] = fmt.Sprintf("//%s%s", g.server.Addr, tpl)
			}
		}
		g.responseJSON(w, http.StatusOK, routes)
	}).Methods("GET").Name("api_directory")

	g.router.HandleFunc("/api/devices", g.getDevices).Methods("GET").Name("devices")
	dev := g.router.PathPrefix("/api/devices").Subrouter()
	// dev.HandleFunc("", g.getDevices).Methods("GET").Name("devices") // doesn't work
	dev.HandleFunc("/{mac}", g.getDevice).Methods("GET").Name("device")
	dev.HandleFunc("/{mac}/upgrade", g.upgradeDevice).Methods("POST").Name("upgrade_device")
	dev.HandleFunc("/{mac}/provision", g.provisionDevice).Methods("POST").Name("provision_device")
	dev.HandleFunc("/{mac}/reboot", g.rebootDevice).Methods("POST").Name("reboot_device")
}

func (g *goWeb) statusJSON(w http.ResponseWriter, status int, message string, v ...interface{}) {
	json := map[string]interface{}{
		"type":    "info",
		"message": fmt.Sprintf(message, v...),
	}

	if 200 <= status && status < 300 {
		json["type"] = "success"
	} else if 400 <= status {
		json["type"] = "danger"
	}

	g.responseJSON(w, status, json)
}

func (g *goWeb) responseJSON(w http.ResponseWriter, status int, v interface{}) {
	if result, err := json.Marshal(v); err != nil {
		log.Printf("[responseJSON error] %v", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
	} else {
		w.Header().Set(headerContentType, contentTypeJSON)
		w.Header().Set(headerContentLength, strconv.Itoa(len(result)))
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	}
}
