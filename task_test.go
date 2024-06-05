package goasynctask

import (
	"fmt"
	"testing"
	"time"
)

func Test(t *testing.T) {
	task := New[string]()
	task.Add(func() (string, error) {
		fmt.Println(1)
		return "1", nil
	})
	task.Add(func() (string, error) {
		fmt.Println(1)
		return "2", nil
	})
	task.Add(func() (string, error) {
		fmt.Println(2)
		return "3", nil
	})
	task.Add(func() (string, error) {
		fmt.Println(2)
		return "3", nil
	})
	task.Add(func() (string, error) {
		fmt.Println(2)
		return "2", nil
	})
	task.Add(func() (string, error) {
		fmt.Println(2)
		return "2", nil
	})
	err := task.Run(time.Second)
	if err != nil {
		t.Log(err)
	}
}
