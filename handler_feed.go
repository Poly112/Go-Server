package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user Users) {

	// Create a context for the request
	ctx := r.Context()

	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v \n", err))
		return
	}

	feed := Feeds{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	}

	result := apiCfg.DB.WithContext(ctx).Create(&feed)
	if result.Error != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create feed: %v \n", result.Error))
		return
	}

	respondWithJSON(w, 201, feed)
}

func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	var feeds []Feeds
	result := apiCfg.DB.Find(&feeds)
	if result.Error != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get feeds: %v \n", result.Error))
		return
	}
	respondWithJSON(w, 200, feeds)
}
