package cache

//var cachePrefix = "crs:micro-goods:"

type Cache interface {
	Key(unique string) string
	MinTime() int
	MaxTime() int
	Set(key string, data interface{})
	Get(key string, targetObject interface{}) interface{}
}

type BasicCache struct {
}

func (b BasicCache) Key(unique string) string {
	return "" + unique
}

func (b BasicCache) MinTime() int {
	return 10
}

func (b BasicCache) MaxTime() int {
	return 100
}

func (b BasicCache) Set(unique string, data interface{}) {
	/*
		redis := db.Redis.Get()
		jsonString, _ := json.Marshal(data)
		expire := rand.Intn(b.MaxTime()) - b.MinTime()
		_ = redis.Send("setex", cachePrefix+b.Key(unique), expire, jsonString)
		_ = redis.Flush()
		_ = redis.Close()

	*/
}

func (b BasicCache) get(unique string, targetObject interface{}) interface{} {
	/*
		redis := db.Redis.Get()
		_ = redis.Send("get", cachePrefix+b.Key(unique))
		jsonString, _ := redis.Receive()
		_ = json.Unmarshal([]byte(jsonString.(string)), targetObject)
		return targetObject

	*/
	return nil
}
