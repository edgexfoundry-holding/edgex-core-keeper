//
// Copyright (C) 2020-2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package redis

// Redis commmands used in this project
// Reference: https://redis.io/commands
const (
	COUNT            = "count"
	MULTI            = "MULTI"
	SET              = "SET"
	GET              = "GET"
	EXISTS           = "EXISTS"
	DEL              = "DEL"
	HLEN             = "HLEN"
	HSCAN            = "HSCAN"
	HSET             = "HSET"
	HGET             = "HGET"
	HGETALL          = "HGETALL"
	HEXISTS          = "HEXISTS"
	HDEL             = "HDEL"
	HSETNX           = "HSETNX"
	MATCH            = "MATCH"
	SADD             = "SADD"
	SREM             = "SREM"
	ZADD             = "ZADD"
	ZREM             = "ZREM"
	EXEC             = "EXEC"
	ZRANGE           = "ZRANGE"
	ZREVRANGE        = "ZREVRANGE"
	MGET             = "MGET"
	ZCARD            = "ZCARD"
	ZCOUNT           = "ZCOUNT"
	UNLINK           = "UNLINK"
	ZRANGEBYSCORE    = "ZRANGEBYSCORE"
	ZREVRANGEBYSCORE = "ZREVRANGEBYSCORE"
	LIMIT            = "LIMIT"
	ZUNIONSTORE      = "ZUNIONSTORE"
	ZINTERSTORE      = "ZINTERSTORE"
	TYPE             = "TYPE"
)

const (
	InfiniteMin     = "-inf"
	InfiniteMax     = "+inf"
	GreaterThanZero = "(0"
	DBKeySeparator  = ":"
	WildCard        = "*"
)

// Redis data types
const (
	Hash   = "hash"
	String = "string"
	None   = "none"
)
