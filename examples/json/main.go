package main

import (
	"encoding/json"
	"fmt"

	events "github.com/guiabolso/events-protocol-go"
)

func main() {
	eventJSON := "{\"id\":\"teste-de-id\",\"flowId\":\"teste-de-id-fluxo\",\"payload\":{\"name\":\"teste\"},\"name\":\"event\",\"version\":\"1\",\"metadata\":{\"createdAt\":\"segunda\"},\"identity\":{\"userId\":11291},\"auth\":{}}"
	// event := json.Unmarshal(eventJSON)

	// IDGenerator := func() string {
	// 	return "this-is-an-id"
	// }
	// session := events.RetrieveEventSession(IDGenerator)
	// session.Import(event)

	event := events.FromJSON(eventJSON)

	var event events.Event
	_ = json.Unmarshal([]byte(eventJSON), &event)

	eventToJSON, _ := json.Marshal(event)
	fmt.Println(string(eventToJSON))
}
