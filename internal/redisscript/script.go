package redisscript

import "github.com/redis/go-redis/v9"

var HSetTimeScript = redis.NewScript(`
  local key = KEYS[1]
  local url = ARGV[1]
  local type = ARGV[2]
  local from = ARGV[3]
  local exp = ARGV[4]
  redis.call('HSET', key,'url', url, 'type', type, 'from', from)
  redis.call('EXPIRE', key, exp)
  return 1
`)
var HSetScript = redis.NewScript(`
  local key = KEYS[1]
  local url = ARGV[1]
  local type = ARGV[2]
  local from = ARGV[3]
  local count = ARGV[4]
  local exp = ARGV[5]
  redis.call('HSET', key,'url', url, 'type', type, 'from', from, 'count', count)
  redis.call('EXPIRE', key, exp)
  return 1
`)
var HGetAllTimeScript = redis.NewScript(`
  local hashKey = KEYS[1]
  local hashInfo = redis.call('HGETALL', hashKey)
  local ttl = redis.call('TTL', hashKey)
  return {hashInfo, ttl}
`)
var HGetUrlTypeTimeScript = redis.NewScript(`
  local hashKey = KEYS[1]
  local url = ARGV[1]
  local type = ARGV[2] 
  local hashInfo = redis.call('HMGET', hashKey, url, type)
  local ttl = redis.call('TTL', hashKey)
  return {hashInfo, ttl}
`)
var IncrTimeScript = redis.NewScript(`
 local key = KEYS[1]
 local expiration = tonumber(ARGV[1])
 local newValue = redis.call('INCR', key)
 redis.call('EXPIRE', key, expiration)
 return newValue
`)
