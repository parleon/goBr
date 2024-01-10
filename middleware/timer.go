package middleware

import (
	"log"
	"net/http"
	"time"
)

func Timer(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		log.Println("Request took: ", time.Since(start))
	}
}