package goasynctask

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type GoroutinePool interface {
	Submit(func()) error
}
type AsyncTask[T any] struct {
	cnt           int
	taskSet       map[string]func() (T, error)
	taskResultSet map[string]taskResult[T]
	pool          GoroutinePool
	timeout       bool
}

// 限时并发
func New[T any]() AsyncTask[T] {
	return AsyncTask[T]{
		taskSet:       map[string]func() (T, error){},
		taskResultSet: map[string]taskResult[T]{},
	}

}
func NewWithPool[T any](pool GoroutinePool) AsyncTask[T] {
	return AsyncTask[T]{
		taskSet:       map[string]func() (T, error){},
		taskResultSet: map[string]taskResult[T]{},
		pool:          pool,
	}

}

func (c *AsyncTask[T]) Add(task func() (T, error)) {
	c.cnt++
	c.taskSet[strconv.Itoa(c.cnt)] = task
}
func (c *AsyncTask[T]) AddWithKey(key string, task func() (T, error)) {
	_, ok := c.taskSet[key]
	if ok {
		panic(fmt.Errorf("duplicated key %s", key))
	}
	c.cnt++
	c.taskSet[key] = task
}

type taskResult[T any] struct {
	Key    string
	Err    error
	Result T
	Cost   int64 //ms
}

func (c *AsyncTask[T]) Result() map[string]taskResult[T] {
	return c.taskResultSet
}
func (c *AsyncTask[T]) Run(timeout time.Duration) error {
	taskLen := len(c.taskSet)
	ch := make(chan taskResult[T], taskLen)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	f := func(key string, task func() (T, error)) {
		now := time.Now()
		result, err := task()
		cost := time.Since(now).Milliseconds()
		if !c.timeout {
			ch <- taskResult[T]{
				Key:    key,
				Err:    err,
				Result: result,
				Cost:   cost,
			}
		}
	}
	if c.pool == nil {
		for key, task := range c.taskSet {
			go f(key, task)
		}
	} else {
		for key, task := range c.taskSet {
			err := c.pool.Submit(func() {
				f(key, task)
			})
			if err != nil {
				return err
			}
		}
	}

	fillRspFunc := func(taskRsp taskResult[T]) (done bool) {
		c.taskResultSet[taskRsp.Key] = taskRsp
		return len(c.taskResultSet) >= taskLen
	}
	for {
		select {
		case <-ctx.Done():
			//极端情况下，在超时的一瞬间，监听的两个channel同时有值时，会随机选择一个case，确保回收返回值

		loop:
			for {
				select {
				case taskRsp := <-ch:
					if fillRspFunc(taskRsp) {
						return nil
					}
				default:
					break loop
				}
			}
			close(ch)
			c.timeout = true
			return errors.New("tasks time out")
		case taskRsp := <-ch:
			if fillRspFunc(taskRsp) {
				return nil
			}
		}
	}
}
