package executor

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	. "github.com/spyzhov/healthy/executor/internal/args"
	"github.com/spyzhov/healthy/step"
	"github.com/spyzhov/safe"
)

type AppRedisArgs struct {
	URL      string               `json:"url"`
	Password string               `json:"password"`
	Cmd      string               `json:"cmd"`
	Args     []interface{}        `json:"args"`
	Require  *AppRedisArgsRequire `json:"require"`
}

func (e *Executor) AppRedis(args *AppRedisArgs) (step.Function, error) {
	scope := "app/redis"
	if err := args.Validate(); err != nil {
		return nil, safe.Wrap(err, scope)
	}

	id := fmt.Sprintf("%s_%s", scope, args.URL)
	e.addSetter(id, func() (interface{}, error) {
		return &redis.Pool{
			MaxIdle: 3,
			Wait:    true,
			Dial: func() (redis.Conn, error) {
				conn, err := redis.DialURL(args.URL)
				if err != nil {
					return nil, err
				}
				if args.Password != "" {
					if _, err := conn.Do("AUTH", args.Password); err != nil {
						safe.Close(conn, "Redis connection")
						return nil, err
					}
				}
				return conn, nil
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if time.Since(t) < time.Minute {
					return nil
				}
				_, err := c.Do("PING")
				return err
			},
		}, nil
	})
	exec := func() (interface{}, error) {
		var (
			pool *redis.Pool
			ok   bool
		)
		if conn, err := e.getConnection(id); err != nil {
			return nil, err
		} else if pool, ok = conn.(*redis.Pool); !ok {
			return nil, fmt.Errorf("redis: connection error")
		}
		conn := pool.Get()
		defer safe.Close(conn, "Redis connection")
		return conn.Do(args.Cmd, args.Args...)
	}
	return func() (*step.Result, error) {
		if reply, err := exec(); err != nil {
			return nil, safe.Wrap(err, scope+": execute")
		} else if err = args.Require.Match(reply); err != nil {
			return nil, err
		}

		return step.NewResultSuccess("OK"), nil
	}, nil
}

func (a *AppRedisArgs) Validate() (err error) {
	if a == nil {
		return fmt.Errorf("required arguments")
	}
	if a.URL == "" {
		return fmt.Errorf("url: require")
	}
	if a.Cmd == "" {
		return fmt.Errorf("cmd: require")
	}
	if err = a.Require.Validate(); err != nil {
		return safe.Wrap(err, "require")
	}
	return
}
