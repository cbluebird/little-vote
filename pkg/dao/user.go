package dao

import (
	"context"
	"errors"
	"log"

	"gorm.io/gorm"

	"little-vote/pkg/database"
	"little-vote/pkg/model"
	"little-vote/pkg/redis"
)

var (
	ctx = context.Background()
)

func GetUserInfo(name string) (*model.User, error) {
	var user model.User
	result := database.DB.Where(
		&model.User{
			Name: name,
		},
	).First(&user)
	return &user, result.Error
}

func GetUserInCache(name string) (int, error) {
	val, err := redis.RedisClient.Get(ctx, name).Int()
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return val, nil
}

func SetUserInCache(name string) error {
	user, err := GetUserInfo(name)
	if err != nil {
		return err
	}
	err = redis.RedisClient.Set(context.Background(), user.Name, user.Count, 0).Err()
	return err
}

func IncrUserInCache(name string) error {
	_, err := redis.RedisClient.Incr(context.Background(), name).Result()
	return err
}

func SyncUser(name string) error {
	count, err := GetUserInCache(name)
	if err != nil {
		return err
	}
	var user model.User
	err = database.DB.Model(&model.User{}).Where(&model.User{Name: name}).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		database.DB.Create(&model.User{Name: name, Count: count})
		return nil
	} else if err != nil {
		return err
	}
	database.DB.Model(&model.User{}).Where(&model.User{Name: name}).Update("count", count)
	return nil
}
