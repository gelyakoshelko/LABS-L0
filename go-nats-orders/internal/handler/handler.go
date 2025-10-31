package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go-nats-orders/internal/cache"
)

type Handler struct {
	cache *cache.Cache
}

func New(c *cache.Cache) *Handler {
	return &Handler{cache: c}
}


func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "static/index.html")
}


func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if id == "" {
		http.Error(w, `{"error":"missing id"}`, http.StatusBadRequest)
		return
	}
	if o, ok := h.cache.Get(id); ok {
		json.NewEncoder(w).Encode(o)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, `{"error":"order not found"}`)
}
