package main

import (
	"fmt"
	"net/http"

	"github.com/Poly112/go-server/auth"
)

type authedHandler func(http.ResponseWriter, *http.Request, Users)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth error: %v\n", err))
			return
		}

		var user Users
		result := apiCfg.DB.WithContext(r.Context()).Where("api_key = ?", apiKey).First(&user)

		if result.Error != nil {
			respondWithError(w, 400, fmt.Sprintf("Error querying database: %v \n", result.Error))
			return
		}

		handler(w, r, user)
	}
}
