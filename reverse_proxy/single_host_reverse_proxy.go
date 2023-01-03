package single_host_reverse_proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
)

func main() {
	logger := log.Default()

	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("TEST-HEADER", "TEST-VALUE")
		fmt.Fprintln(w, "this call was relayed by the reverse proxy")
	}))
	defer backendServer.Close()

	rpURL, err := url.Parse(backendServer.URL)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("temporary server started successfully on =>", rpURL)

	proxyServer := httputil.NewSingleHostReverseProxy(rpURL)

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		logger.Println("got an request from =>", r.RemoteAddr)
		proxyServer.ServeHTTP(rw, r)
		logger.Println("proxy answered to =>", r.RemoteAddr)
	})

	http.ListenAndServe(":9090", nil)
}
