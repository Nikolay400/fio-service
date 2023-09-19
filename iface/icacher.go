package iface

import "time"

type Icacher interface {
	Set(key string, value interface{}, t time.Duration) error
	Del(key string) error
	Get(key string) (string, error)
}
