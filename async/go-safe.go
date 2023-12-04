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

func Go2[T1 any, T2 any](fn func(T1, T2)) func(T1, T2) {
	return func(p1 T1, p2 T2) {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Errorm("panic", "err", err)
				}
			}()

			fn(p1, p2)
		}()
	}
}

func Go3[T1 any, T2 any, T3 any](fn func(T1, T2, T3)) func(T1, T2, T3) {
	return func(p1 T1, p2 T2, p3 T3) {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Errorm("panic", "err", err)
				}
			}()

			fn(p1, p2, p3)
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
