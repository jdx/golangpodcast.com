package main

import (
	"log"

	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

var conn *s3.S3
var episodes []*Episode

type Episode struct {
	Title string
	Url   string
	s3.Key
}

func ParseEpisode(key s3.Key) *Episode {
	url := "http://d31dj2i8pqzrq5.cloudfront.net/" + key.Key
	title := key.Key
	return &Episode{title, url, key}
}

func GetEpisodes() []*Episode {
	if episodes != nil {
		return episodes
	}
	log.Println("Fetching episodes...")
	bucket := getBucket()
	resp, err := bucket.List("episodes/20", "", "", 1000)
	if err != nil {
		log.Fatal(err)
	}
	episodes = make([]*Episode, 0, len(resp.Contents))
	for _, key := range resp.Contents {
		episodes = append(episodes, ParseEpisode(key))
	}
	return episodes
}

func GetEpisode(slug string) *Episode {
	key := "episodes/" + slug
	for _, episode := range GetEpisodes() {
		if episode.Key.Key == key {
			return episode
		}
	}
	return nil
}

func getAwsConnection() *s3.S3 {
	if conn != nil {
		return conn
	}
	log.Println("Connecting to S3...")
	auth, err := aws.EnvAuth()
	if err != nil {
		log.Fatal(err)
	}
	conn = s3.New(auth, aws.USEast)
	return conn
}

func getBucket() *s3.Bucket {
	return getAwsConnection().Bucket("golangpodcast")
}
