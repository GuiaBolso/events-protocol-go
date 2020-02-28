package main

import (
	"fmt"

	guuid "github.com/google/uuid"
	events "github.com/guiabolso/events-protocol-go"
)

func main() {
	UUIDGenerator := func() string {
		uuid := guuid.New()
		return uuid.String()
	}
	session := events.GenerateEventSession(UUIDGenerator)

	metadata := map[string]interface{}{
		"criadoEm": "exemplo",
	}
	session.SetMetadata(metadata)

	payload := map[string]interface{}{
		"empresa": "guiabolso",
		"missao":  "melhorar a vida do brasileiro e transformar o sistema financeiro",
	}
	eventTemplate := session.RegisterEvent("uuid:event", "1").WithPayload(payload)

	event1 := eventTemplate.Prepare()

	fmt.Print(event1)

	event2 := eventTemplate.Prepare()

	fmt.Print(event2)
}
