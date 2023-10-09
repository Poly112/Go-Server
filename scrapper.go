package main

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func startScrapping(db *gorm.DB, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scraping on %v goroutines every %v duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)

	for ; ; <-ticker.C {
		var feeds []Feeds

		result := db.WithContext(context.Background()).Order("`last_fetched_at` ASC").Limit(concurrency).Find(&feeds)

		if result.Error != nil {
			log.Println("error fetching feeds:", result.Error)
			continue
		}

		wg := &sync.WaitGroup{}

		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()

	}
}

func scrapeFeed(db *gorm.DB, wg *sync.WaitGroup, feed Feeds) {
	defer wg.Done()

	result := db.WithContext(context.Background()).Model(&feed).
		Where("id = ?", feed.ID).
		Updates(map[string]interface{}{"last_fetched_at": time.Now(), "updated_at": time.Now()})

	if result.Error != nil {
		log.Println("error marking feed as fetched:", result.Error)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Panicln("Error fetching feed:", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Panicf("Couldn't parse date %v with err %v", item.PubDate, err)
			continue
		}
		post := Posts{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Description: &item.Description,
			PublishedAt: pubAt,
			Url:         item.Link,
			FeedID:      feed.ID,
		}
		result := db.WithContext(context.Background()).Create(&post)
		if result.Error != nil {
			if strings.Contains(result.Error.Error(), "Duplicate") {
				continue
			}
			log.Panicln("Failed to create post:", result.Error)
			return
		}
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
