package database

import (
	"context"
	"fmt"
	"github.com/musicorum-app/resource-manager/structs"
	"github.com/musicorum-app/resource-manager/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

var client *mongo.Client

func Initialize() {
	var err error
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(utils.GetEnvVar("MONGODB")))
	if err != nil {
		println("Could not connect to database. Please check your access")
	}
	utils.FailOnError(err)

	err = client.Ping(ctx, readpref.Primary())
	utils.FailOnError(err)
}

func FindArtist(resource string, key string) *structs.ArtistCache {
	hash := utils.Hash(key)
	collection := client.Database("resources").Collection(resource)
	fmt.Println("DOING MONGO SEARCH: " + key)
	filter := bson.D{{"hash", hash}}
	var result *structs.ArtistCache
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil
	}
	return result
}

func InsertArtist(artist string, spotify string, url string) {
	collection := client.Database("resources").Collection("artists")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := collection.InsertOne(ctx, bson.M{
		"key":      artist,
		"spotify":  spotify,
		"hash":     utils.Hash(artist),
		"url":      url,
		"cachedAt": time.Now().Unix() * 1000,
	})
	if err != nil {
		println("ERROR WHILE SAVING ON DATABASE")
		println(err)
	}
}

func FindAlbum(album string, artist string) *structs.AlbumCache {
	fmt.Println("Starting to search on database")

	hash := utils.HashAlbum(album, artist)
	collection := client.Database("resources").Collection("albums")
	filter := bson.D{{"hash", hash}}
	var result *structs.AlbumCache
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil
	}
	return result
}

func InsertAlbum(cache *structs.AlbumResponse) {
	collection := client.Database("resources").Collection("albums")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := collection.InsertOne(ctx, bson.M{
		"hash":     utils.HashAlbum(cache.Name, cache.Artist),
		"name":     cache.Name,
		"artist":   cache.Artist,
		"spotify":  cache.Spotify,
		"cover":    cache.Cover,
		"cachedAt": time.Now().Unix() * 1000,
	})
	if err != nil {
		println("ERROR WHILE SAVING ON DATABASE")
		println(err)
	}
}

func FindTrack(hash string) *structs.TrackCache {
	fmt.Println("Starting to search on database")

	collection := client.Database("resources").Collection("tracks")
	filter := bson.D{{"hash", hash}}
	var result *structs.TrackCache
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil
	}
	return result
}

func SetTrack(cache *structs.TrackResponse) {
	collection := client.Database("resources").Collection("tracks")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := collection.InsertOne(ctx, bson.M{
		"hash":     utils.HashTrack(cache.Name, cache.Artist, cache.Album),
		"name":     cache.Name,
		"artist":   cache.Artist,
		"album":    cache.Album,
		"cover":    cache.Cover,
		"spotify":  cache.Spotify,
		"duration": cache.Duration,
		"cachedAt": time.Now().Unix() * 1000,
	})
	if err != nil {
		println("ERROR WHILE SAVING ON DATABASE")
		println(err)
	}
}
