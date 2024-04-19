package main

import (
	"log"
	"net/http"

	"github.com/Friday1602/chirpy/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
	db             *database.DB
}
type chripyParams struct {
	Body string `json:"body"`
}
type user struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CustomClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	mux := http.NewServeMux()
	apiCfg := &apiConfig{}
	fileServer := http.FileServer(http.Dir("./app"))
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer))) //* for wildcard

	apiCfg.db, err = database.NewUserDB("userDatabase.json")
	if err != nil {
		log.Fatal(err)
	}

	mux.HandleFunc("GET /admin/metrics", apiCfg.metrics)

	mux.HandleFunc("/api/reset", apiCfg.reset)

	fileServer = http.FileServer(http.Dir("./app/assets"))
	mux.Handle("/app/assets/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/assets", fileServer)))

	mux.HandleFunc("GET /api/healthz", readiness)

	mux.HandleFunc("POST /api/chirps", apiCfg.validateChirpy)
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirpy)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirpyFromID)
	mux.HandleFunc("POST /api/users", apiCfg.createUser)
	mux.HandleFunc("POST /api/login", apiCfg.userValidation)
	mux.HandleFunc("PUT /api/users", apiCfg.updateUser)
	mux.HandleFunc("POST /api/refresh", apiCfg.refreshTokenAuth)
	mux.HandleFunc("POST /api/revoke", apiCfg.revokeToken)

	corsMux := middlewareCors(mux)
	log.Print("starting server on :8080")
	err = http.ListenAndServe(":8080", corsMux)
	log.Fatal(err)
}
