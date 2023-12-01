package async

import "github.com/night-sword/kratos-kit/log"

func Go[T any](fn func(T)) func(T) {
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

func Run(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Errorm("panic", "err", err)
			}
		}()

		fn()
	}()
}
