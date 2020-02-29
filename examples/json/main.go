package main

import (
	"fmt"

	events "github.com/guiabolso/events-protocol-go"
)

func main() {
	eventJSON := "{\"id\":\"teste-de-id\",\"flowId\":\"teste-de-id-fluxo\",\"payload\":{\"name\":\"teste\"},\"name\":\"event\",\"version\":\"1\",\"metadata\":{\"createdAt\":\"segunda\"},\"identity\":{\"userId\":11291},\"auth\":{}}"
	// event := json.Unmarshal(eventJSON)

	session, _ := events.ImportJSONEventSession(eventJSON)

	fmt.Println(session.Events)

	ok, registeredEvent := session.Events["event"]

	fmt.Println(ok)
	fmt.Println(registeredEvent)

	eventTemplate := session.RegisterEvent("teste", "1")

	event := eventTemplate.Prepare()

	fmt.Println(event.ToJSON())
}
