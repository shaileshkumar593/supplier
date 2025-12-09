package common

type RedisKey struct {
	Key string `json:"key" binding:"required" validate:"required"`
}

type PluToRedis struct {
	Plu map[string]string `json:"plu" binding:"required" validate:"required"`
}
