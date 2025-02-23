package main

import (
    "backend/models"
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "time"

    "github.com/pascaldekloe/jwt"
    "golang.org/x/crypto/bcrypt"
)



type Credentials struct {
    Username string 
    Password string 
}

func (app *application) signIn(w http.ResponseWriter, r *http.Request) {
    var creds Credentials

    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        app.errorJSON(w, errors.New("invalid request payload"))
        return
    }

    user, err := app.models.DB.GetUserByUsername(creds.Username)
    if err != nil {
        app.errorJSON(w, errors.New("invalid username or password"))
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
    if err != nil {
        app.errorJSON(w, errors.New("invalid username or password"))
        return
    }

    var claims jwt.Claims
    claims.Subject = fmt.Sprint(user.UserID)
    claims.Issued = jwt.NewNumericTime(time.Now())
    claims.NotBefore = jwt.NewNumericTime(time.Now())
    claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
    claims.Issuer = "testing"
    claims.Audiences = []string{"testing"}

    jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secret))
    if err != nil {
        app.errorJSON(w, errors.New("error signing tok"))
        return
    }

    response := map[string]interface{}{
        "token":    string(jwtBytes), // Convert jwtBytes to string
        "user_id":  user.UserID,
        "username": user.Username,
    }

    app.writeJSON(w, http.StatusOK, response, "response")
}

func (app *application) signUp(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        app.errorJSON(w, errors.New("invalid request payload"))
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 12)
    if err != nil {
        app.errorJSON(w, errors.New("error hashing password"))
        return
    }

    user := &models.User{
        Username: creds.Username,
        Password: string(hashedPassword),
    }

    err = app.models.DB.InsertUser(user)
    if err != nil {
        app.errorJSON(w, errors.New("error inserting user"))
        return
    }

    response := map[string]interface{}{
        "user_id":  user.UserID,
        "username": user.Username,
    }

    app.writeJSON(w, http.StatusCreated, response, "response")
}


