package employee

import (
	"fmt"
	"strings"

	"github.com/garyburd/redigo/redis"
)

func SearchName(pool *redis.Pool, name string) ([]string, error) {
	var (
		keys         []string
		matchedNames []string
	)

	conn := pool.Get()
	defer conn.Close()

	cursor := 0
	pattern := fmt.Sprintf("hr:employees:index:name_to_ids:%v*", name)

	for {
		v, err := redis.Values(conn.Do("SCAN", cursor, "MATCH", pattern, "COUNT", 1024))
		if err != nil {
			return []string{}, err
		}

		if v, err = redis.Scan(v, &cursor, &keys); err != nil {
			return []string{}, err
		}

		for _, k := range keys {
			matchedNames = append(matchedNames, strings.TrimPrefix(k, "hr:employees:index:name_to_ids:"))
		}

		if cursor == 0 {
			break
		}
	}
	return matchedNames, nil
}
