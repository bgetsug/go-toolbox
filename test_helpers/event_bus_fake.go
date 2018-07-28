package test_helpers

import (
	"github.com/asaskevich/EventBus"
	"github.com/bgetsug/go-toolbox/logging"
	"go.uber.org/zap"
)

type fakeEventBus struct {
	log *zap.SugaredLogger
}

// NewFakeEventBus returns a pointer to an EventBus.Bus that does not execute any event handlers.
// Instead, it logs debug messages for published events, and does so synchronously.
func NewFakeEventBus() EventBus.Bus {
	return &fakeEventBus{logging.NewModuleLog("test", "fake_event_bus")}
}

func (b *fakeEventBus) HasCallback(topic string) bool {
	b.log.Debugf("HasCallback(%s)", topic)
	return false
}

func (b *fakeEventBus) Publish(topic string, args ...interface{}) {
	b.log.Debugf("Publish(%s, %s)", topic, args)
}

func (b *fakeEventBus) Subscribe(topic string, fn interface{}) error {
	b.log.Debugf("Subscribe(%s)", topic)
	return nil
}

func (b *fakeEventBus) SubscribeAsync(topic string, fn interface{}, transactional bool) error {
	b.log.Debugf("SubscribeAsync(%s)", topic)
	return nil
}

func (b *fakeEventBus) SubscribeOnce(topic string, fn interface{}) error {
	b.log.Debugf("SubscribeOnce(%s)", topic)
	return nil
}

func (b *fakeEventBus) SubscribeOnceAsync(topic string, fn interface{}) error {
	b.log.Debugf("SubscribeOnceAsync(%s)", topic)
	return nil
}

func (b *fakeEventBus) Unsubscribe(topic string, handler interface{}) error {
	b.log.Debugf("Unsubscribe(%s)", topic)
	return nil
}

func (b *fakeEventBus) WaitAsync() {
	b.log.Debug("WaitAsync()")
}
