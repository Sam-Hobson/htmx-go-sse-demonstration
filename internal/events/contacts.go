package events

import (
	"errors"
	"sync"
	"time"
)

type addContactQueue struct {
	queueTimeout  time.Duration
	subscriptions *sync.Map
}

var ContactQueue = &addContactQueue{
	queueTimeout:  2 * time.Second,
	subscriptions: &sync.Map{},
}

func key(k any) string {
	return k.(string)
}

func value(v any) chan *Contact {
	return v.(chan *Contact)
}

func (c *addContactQueue) Queue(data Contact) {
	addDataToAllChannels(c.queueTimeout, c.subscriptions, &data)
}

func (c *addContactQueue) Subscribe(key string, chanSize int) <-chan *Contact {
	ch := make(chan *Contact, chanSize)
	c.subscriptions.Store(key, ch)
	return ch
}

func (c *addContactQueue) Unsubscribe(key string) {
	ch, loaded := c.subscriptions.LoadAndDelete(key)
	if loaded {
		close(value(ch))
	}
}

func (c *addContactQueue) Close() {
	closeAll[chan *Contact](c.subscriptions)
}

type removeContactQueue struct {
	queueTimeout  time.Duration
	subscriptions *sync.Map
}

var RemoveContactQueue = &removeContactQueue{
	queueTimeout:  2 * time.Second,
	subscriptions: &sync.Map{},
}

func (c *removeContactQueue) Queue(data Contact) {
	addDataToAllChannels(c.queueTimeout, c.subscriptions, &data)
}

func (c *removeContactQueue) Subscribe(key string, chanSize int) <-chan *Contact {
	ch := make(chan *Contact, chanSize)
	c.subscriptions.Store(key, ch)
	return ch
}

func (c *removeContactQueue) Unsubscribe(key string) {
	ch, loaded := c.subscriptions.LoadAndDelete(key)
	if loaded {
		close(value(ch))
	}
}

func (c *removeContactQueue) Close() {
	closeAll[chan *Contact](c.subscriptions)
}

var contactId = 0

type Contact struct {
	Id        int
	FirstName string  `json:"firstName" form:"firstName"`
	LastName  string  `json:"lastName" form:"lastName"`
	Age       int     `json:"age" form:"age"`
	Height    float32 `json:"height" form:"height"`
	Gender    string  `json:"gender" form:"gender"`
}

func (c Contact) Validate() error {
	if c.FirstName == "" {
		return errors.New("Invalid first name")
	}
	if c.LastName == "" {
		return errors.New("Invalid last name")
	}
	if c.Age < 0 {
		return errors.New("Invalid age")
	}
	if c.Height <= 0 {
		return errors.New("Invalid height")
	}
	if c.Gender != "M" && c.Gender != "F" {
		return errors.New("Invalid gender")
	}

	return nil
}

func NewContact() *Contact {
	contactId++
	return &Contact{
		Id: contactId,
	}
}
