package main

import (
	"encoding/json"
	"fmt"
	"main/config"
	"main/route"
	"net/http"
	"os"

	"github.com/ReneKroon/ttlcache"
	"github.com/gorilla/mux"
)

func main() {
	// 1.Intialize DataBase Connection
	db, err := config.Connect()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		os.Exit(1)
	}
	defer db.Close()

	//2. Initialize router
	router := mux.NewRouter()

	//3. Define routes
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Hello, World!", "success": "true"})
	})

	//for caching our otp for email verification
	cache := ttlcache.NewCache()
	defer cache.Close()
	r := route.Route{Db: db, Router: router, Cache: cache}
	r.Routes()

	//4.start server
	fmt.Println("Starting server on port 8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}

}
