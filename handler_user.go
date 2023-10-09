package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	// Create a context for the request
	ctx := r.Context()

	type parameters struct {
		Name string `json:"name"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v \n", err))
		return
	}

	user := Users{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	}

	result := apiCfg.DB.WithContext(ctx).Create(&user)
	if result.Error != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create User: %v \n", result.Error))
		return
	}

	respondWithJSON(w, 201, user)
}

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user Users) {

	respondWithJSON(w, 200, user)

}
func (apiCfg *apiConfig) handlerGetPostsForUser(w http.ResponseWriter, r *http.Request, user Users) {
	var posts []Posts

	result := apiCfg.DB.WithContext(r.Context()).Table("posts").
		Select("posts.*").
		Joins("join feed_follows on posts.feed_id = feed_follows.feed_id").
		Where("feed_follows.user_id = ?", user.ID).
		Order("posts.published_at desc").
		Limit(10).
		Find(&posts)

	if result.Error != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get posts: %v\n", result.Error))
		return
	}
	respondWithJSON(w, 200, posts)
}
