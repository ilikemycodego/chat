package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gorilla/mux"
)

func RobiProxy(r *mux.Router) {
	robiAddr := os.Getenv("ROBI_URL")
	if robiAddr == "" {
		robiAddr = "http://localhost:8083"
	}

	robiURL, err := url.Parse(robiAddr)
	if err != nil {
		log.Println("Неверный адрес robi:", err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(robiURL)

	r.PathPrefix("/robi").Handler(http.StripPrefix("/robi", proxy))

	log.Println("Проксируем /robi →", robiAddr)
}
