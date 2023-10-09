package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeedFollows(w http.ResponseWriter, r *http.Request, user Users) {

	// Create a context for the request
	ctx := r.Context()

	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v \n", err))
		return
	}

	feedFollows := FeedFollows{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	}

	result := apiCfg.DB.WithContext(ctx).Create(&feedFollows)
	if result.Error != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create feed: %v \n", result.Error))
		return
	}

	respondWithJSON(w, 201, feedFollows)
}

func (apiCfg *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user Users) {

	// Create a context for the request
	ctx := r.Context()

	var feedFollows []FeedFollows
	result := apiCfg.DB.WithContext(ctx).Where("user_id", user.ID).Find(&feedFollows)
	if result.Error != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't find feed follows: %v \n", result.Error))
		return
	}

	respondWithJSON(w, 200, feedFollows)
}

func (apiCfg *apiConfig) handlerDeleteFeedFollows(w http.ResponseWriter, r *http.Request, user Users) {

	// Create a context for the request
	ctx := r.Context()
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	if feedFollowIDStr == "" {
		respondWithError(w, 400, fmt.Sprintf("No feed follow id found: %v \n", feedFollowIDStr))
		return
	}
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't parse feed follow id: %v \n", err))
		return
	}

	var feedFollows FeedFollows
	result := apiCfg.DB.WithContext(ctx).Where("user_id = ? AND id = ?", user.ID, feedFollowID).Delete(&feedFollows)
	if result.Error != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't find feed follows: %v \n", result.Error))
		return
	}

	respondWithJSON(w, 200, struct{}{})
}
