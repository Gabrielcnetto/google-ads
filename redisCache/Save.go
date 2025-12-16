package rediscache

import (
	"context"
	"encoding/json"
	"netto/models"
	"time"
)

func SaveAccountsOnCache(data []models.GoogleAdwAccount, ctx context.Context) error {
	redisconn := RedisConnection()
	fmtData, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	return redisconn.Set(ctx, "adw_accounts", fmtData, time.Duration(time.Second*30)).Err()
}
