package internal

import "log/slog"

type Subscriber chan string

var (
	register   = make(chan Subscriber)
	unregister = make(chan Subscriber)
	publish    = make(chan string)
)

func StartActor() {
	go func() {
		subscribers := make(map[Subscriber]struct{})

		for {
			select {

			// New subscriber joins
			case sub := <-register:
				subscribers[sub] = struct{}{}

			// Subscriber leaves
			case sub := <-unregister:
				if _, ok := subscribers[sub]; ok {
					delete(subscribers, sub)
					close(sub)
				}

			// New message published to channel
			case msg := <-publish:
				slog.Info("Actor publishing to channel", "msg", msg)
				for sub := range subscribers {
					select {
					case sub <- msg:
						slog.Info("Message delivered to subscriber")
					default:
						slog.Info("Buffer full")
						sub <- msg
					}

				}
			}
		}
	}()
}

// Subscribe returns a buffered channel that receives new messages.
func Subscribe() Subscriber {
	sub := make(Subscriber, 10)
	register <- sub
	return sub
}

// Unsubscribe removes and closes the subscriber channel.
func Unsubscribe(sub Subscriber) {
	unregister <- sub
}

// Publish sends a message to the actor
func Publish(msg string) {
	publish <- msg
}
