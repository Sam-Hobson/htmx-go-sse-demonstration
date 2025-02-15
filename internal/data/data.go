package data

import (
	"salad2/internal/events"
	"salad2/internal/utils"
)

type MainDataStruct struct {
	ContactList *utils.ConcurrentSlice[*events.Contact]
}

var MainData *MainDataStruct = &MainDataStruct{
	ContactList: utils.NewConcurrentSlice[*events.Contact](),
}

func init() {
	go func() {
		for {
			select {
			case data := <-events.ContactQueue.Subscribe("MainData", 100):
				MainData.ContactList.Append(data)
			}
		}
	}()

	go func() {
		for {
			select {
			case data := <-events.RemoveContactQueue.Subscribe("MainData", 100):
				MainData.ContactList.DeleteFunc(func(c *events.Contact) bool { return c.Id == data.Id })
			}
		}
	}()
}
