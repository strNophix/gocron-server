package gocron_server_test

import "fmt"

type Counter struct {
	Current int
}

func (c *Counter) Increment() (string, error) {
	c.Current += 1
	return fmt.Sprint(c.Current), nil
}
