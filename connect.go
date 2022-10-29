package main

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func Client() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return rdb
}

func addData(data string, key string) error {
	client := Client()
	var dat []string

	val, err := client.Get(ctx, key).Bytes()
	if err != nil {
		j, e := json.Marshal([]string{data})
		err := client.Set(ctx, key, j, 0).Err()
		if err != nil || e != nil {
			return err
		}
	}

	if len(val) > 0 {
		val, err := client.Get(ctx, key).Bytes()
		if err != nil {
			return err
		}
		e := json.Unmarshal(val, &dat)
		if e != nil {
			return err
		}

		dat = append(dat, data)
		j, e := json.Marshal(&dat)
		er := client.Set(ctx, key, j, 0).Err()
		if er != nil || e != nil {
			return err
		}

	}
	return nil

}

func getData(key string) ([]string, error) {
	client := Client()
	var dat = []string{}

	data, err := client.Get(ctx, key).Bytes()
	if err != nil {
		return dat, err
	}
	er := json.Unmarshal(data, &dat)
	if er != nil {
		return dat, er
	}
	return dat, nil

}
