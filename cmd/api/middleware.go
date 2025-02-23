package main

import (
    "net/http"
    "strings"
    "errors"
    "github.com/pascaldekloe/jwt"
    "time"	
    "strconv"
    "log"
)

func (app *application) enableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

        if r.Method == http.MethodOptions {
            return
        }

        next.ServeHTTP(w, r)
    })
}

func (app *application) checkToken(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Vary", "Auhtorization")

        authHeader := r.Header.Get("Authorization")

        if authHeader == "" {

        }

        headerParts := strings.Split(authHeader, " ")
        if len(headerParts) != 2 {
            app.errorJSON(w, errors.New("invalid auth header"))
            return
        }

        if headerParts[0] != "Bearer" {
            app.errorJSON(w, errors.New("unauthorized - no bearer"))
            return
        }
        
        token := headerParts[1]

        claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.jwt.secret))
        if err != nil {
            app.errorJSON(w, errors.New("unauthorized - failed hmac check"), http.StatusForbidden)
            return
        }

        if !claims.Valid(time.Now()) {
            app.errorJSON(w, errors.New("unauthorized - failed hmac check"), http.StatusForbidden)
            return
        }

        if !claims.AcceptAudience("testing") {
            app.errorJSON(w, errors.New("invalid audience"), http.StatusForbidden)
            return
        }

        if claims.Issuer != "testing" {
            app.errorJSON(w, errors.New("invalid issuer"), http.StatusForbidden)
            return
        }

        userID, err := strconv.ParseInt(claims.Subject,10,64)
        if err != nil {
            app.errorJSON(w, errors.New("invalid subject"), http.StatusForbidden)
            return 
        }

        log.Println("Valid user:", userID)




        next.ServeHTTP(w, r)
    })
}