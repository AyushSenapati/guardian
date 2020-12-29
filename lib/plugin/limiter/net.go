package limiter

import (
	"log"
	"net"
	"net/http"
	"strings"
)

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
		if ip != "" {
			ip = strings.SplitN(ip, ",", 2)[0]
		}
	}

	if ip == "" {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Printf("error: request IP: %q is not IP:Port", r.RemoteAddr)
		}
		ip = host
	}

	ip = net.ParseIP(ip).String()
	if ip == "" {
		log.Printf("error: %q is not a valid IP", ip)
	}

	return ip
}
