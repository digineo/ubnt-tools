package web

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/digineo/ubnt-tools/provisioner/ui"
)

// GET /
func (g *goWeb) getRoot(w http.ResponseWriter, r *http.Request) {
	asset, err := ui.Asset("views/index.html")

	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
		return
	}

	w.Header().Set(headerContentType, contentTypeHTML)
	w.Header().Set(headerContentLength, strconv.Itoa(len(asset)))
	w.WriteHeader(http.StatusOK)
	w.Write(asset)
}

// GET /assets/*file
func (g *goWeb) getStaticAsset(w http.ResponseWriter, r *http.Request) {
	assetPath := r.URL.Path[1:] // remove leading slash
	asset, err := ui.Asset(assetPath)

	if err != nil {
		log.Print(err)
		http.NotFound(w, r)
		return
	}

	i, err := ui.AssetInfo(assetPath)
	if err != nil {
		log.Printf("Error reading static asset %s: %v", assetPath, err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	modTime := i.ModTime()
	contentType := "application/octed-stream"
	if ct := mime.TypeByExtension(filepath.Ext(assetPath)); ct != "" {
		contentType = ct
	}

	t, err := time.Parse(httpDateFormat, r.Header.Get(headerIfModifiedSince))
	if err == nil && modTime.Before(t.Add(cacheDuration)) {
		w.Header().Del(headerContentType)
		w.Header().Del(headerContentLength)
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.Header().Set(headerContentType, contentType)
	w.Header().Set(headerLastModified, modTime.UTC().Format(httpDateFormat))
	w.Header().Set(headerContentLength, fmt.Sprintf("%d", len(asset)))
	w.WriteHeader(http.StatusOK)
	w.Write(asset)
}
