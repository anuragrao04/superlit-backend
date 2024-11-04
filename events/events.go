package events

import "github.com/asaskevich/EventBus"

var bus = EventBus.New()

func PublishUserCreated(email string) {
	bus.Publish("user:created", email)
}

func SubscribeUserCreated(handler func(string)) {
	bus.Subscribe("user:created", handler)
}
