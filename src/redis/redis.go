package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/musicorum-app/resource-manager/structs"
	"github.com/musicorum-app/resource-manager/utils"
	"time"
)

var ctx = context.Background()
var rdb *redis.Client
var duration time.Duration

func InitializeRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     utils.GetEnvVar("REDIS_ADDR"),
		Password: utils.GetEnvVar("REDIS_PASS"),
		DB:       0,
	})

	var err error
	duration, err = time.ParseDuration(utils.GetEnvVar("REDIS_EXPIRY"))
	if err != nil {
		println("Duration not valid, using default.")
		duration, _ = time.ParseDuration("100h")
	}

	pong, err := rdb.Ping(ctx).Result()
	fmt.Println(pong, err)
}

func SetArtist(artist *structs.ArtistResponse) {
	_, errG := rdb.Get(ctx, artist.Hash).Result()
	if errG != redis.Nil {
		return
	}

	fmt.Println("SETTING ARTIST ON REDIS")
	jsonData, err := json.Marshal(artist)
	fmt.Println("JSON " + string(jsonData))
	if err != nil {
		println("REDIS ERROR")
		println(err.Error())
	}

	err = rdb.Set(ctx, artist.Hash, string(jsonData), duration).Err()
	if err != nil {
		println("REDIS ERROR")
		println(err.Error())
	}
}

func FindArtist(hash string) *structs.ArtistResponse {
	fmt.Println("Searching on redis for " + hash)
	result, err := rdb.Get(ctx, hash).Result()
	if err != nil {
		return nil
	}
	var data *structs.ArtistResponse
	err = json.Unmarshal([]byte(result), &data)
	if err != nil {
		fmt.Println("Error while parsing json from Redis")
		return nil
	}
	return data
}

func SetAlbum(album *structs.AlbumResponse) {
	_, errG := rdb.Get(ctx, album.Hash).Result()
	if errG != redis.Nil {
		return
	}

	fmt.Println("SETTING ALBUM ON REDIS")
	jsonData, err := json.Marshal(album)
	fmt.Println("JSON " + string(jsonData))
	if err != nil {
		println(err.Error())
	}

	err = rdb.Set(ctx, album.Hash, string(jsonData), duration).Err()
	if err != nil {
		println(err.Error())
	}
}

func FindAlbum(hash string) *structs.AlbumResponse {
	fmt.Println("Searching on redis for album " + hash)
	result, err := rdb.Get(ctx, hash).Result()
	if err != nil {
		return nil
	}
	var data *structs.AlbumResponse
	err = json.Unmarshal([]byte(result), &data)
	if err != nil {
		fmt.Println("Error while parsing json from Redis")
		return nil
	}
	return data
}

func SetTrack(track *structs.TrackResponse) {
	_, errG := rdb.Get(ctx, track.Hash).Result()
	if errG != redis.Nil {
		return
	}

	jsonData, err := json.Marshal(track)

	if err != nil {
		println(err.Error())
	}

	err = rdb.Set(ctx, track.Hash, string(jsonData), duration).Err()
	if err != nil {
		println(err.Error())
	}
}

func FindTrack(hash string) *structs.TrackResponse {
	result, err := rdb.Get(ctx, hash).Result()
	if err != nil {
		utils.FailOnError(err)
		return nil
	}

	var data *structs.TrackResponse
	err = json.Unmarshal([]byte(result), &data)
	if err != nil {
		fmt.Println("Error while parsing json from Redis")
		utils.FailOnError(err)
		return nil
	}
	return data

}
