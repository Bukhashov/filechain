package utils

import "time"

func DoWithTries(fn func() error, attems int, dilay time.Duration)(err error){
	for attems > 0 {
		if err = fn(); err != nil {
			time.Sleep(dilay)
			attems--
			continue
		}
		return nil
	}
	return
}