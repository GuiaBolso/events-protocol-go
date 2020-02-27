package main

import (
	"encoding/json"
	"fmt"
	guuid "github.com/google/uuid"
	events "github.com/guiabolso/events-protocol-go"
)

func main() {
	UUIDGenerator := func() string {
		uuid := guuid.New()
		return uuid.String()
	}
	session := events.RetrieveEventSession(UUIDGenerator)

	event := session.RegisterEvent("json:parseable:event", "1")

	fmt.PrintLn(json.Marshal(event))
}
