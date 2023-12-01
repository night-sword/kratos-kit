package util

import "github.com/night-sword/kratos-kit/log"

func GoSafe[T any](fn func(T)) func(T) {
	return func(p T) {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Errorm("panic", "err", err)
				}
			}()

			fn(p)
		}()
	}
}

func GoSafeRun(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Errorm("panic", "err", err)
			}
		}()

		fn()
	}()
}
