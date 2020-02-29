package events

import (
	"fmt"
	"math/rand"
	"testing"
)

func Test_Session(t *testing.T) {
	var basicGenerator SessionIDGenerator
	count := 0
	basicGenerator = func() string {
		count++
		return fmt.Sprintf("%d", count)
	}

	session := GenerateEventSession(basicGenerator)

	if session.SessionID != "1" {
		t.Error("Expecting sessionId to be 1:", session.SessionID)
	}

	session2 := GenerateEventSession(basicGenerator)

	if session2.SessionID == session.SessionID {
		t.Error("Expecting sessions to be different:", session2.SessionID, session.SessionID)
	}
	count = 0
}

func Test_Event(t *testing.T) {
	eventName := "event:name"
	mockIDGenerator := func() string {
		return fmt.Sprintf("%d", rand.Int())
	}
	session := GenerateEventSession(mockIDGenerator)

	identity := map[string]interface{}{
		"identity": "identity",
	}
	session.SetIdentity(identity)

	metadata := map[string]interface{}{
		"metadata": "metadata",
	}
	session.SetMetadata(metadata)

	templateEvent := session.RegisterEvent(eventName, "1")

	if len(session.Events) != 1 {
		t.Error("Expecting event list to be 1", session)
	}

	payload := map[string]interface{}{
		"payload": "payload",
	}
	templateEvent.WithPayload(payload)

	event1 := templateEvent.Prepare()

	if event1.Identity["identity"] != identity["identity"] {
		t.Error("Expecting template identity to be used in final event", event1)
	}
	if event1.Metadata["metadata"] != metadata["metadata"] {
		t.Error("Expecting template metadata to be used in final event", event1)
	}
	if event1.Payload["payload"] != payload["payload"] {
		t.Error("Expecting template payload to be used in final event", event1)
	}

	event2 := templateEvent.Prepare()

	if event1.ID != event2.ID {
		t.Error("Expecting multiple event with same session to have same ID", event1.ID, event2.ID)
	}
	if event1.FlowID == event2.FlowID {
		t.Error("Expecting multiple event not have same FlowID", event1.FlowID, event2.FlowID)
	}
}

func Test_Event_JSON(t *testing.T) {
	eventName := "event:name"
	mockIDGenerator := func() string {
		return "teste-de-id"
	}
	session := GenerateEventSession(mockIDGenerator)

	templateEvent := session.RegisterEvent(eventName, "1")

	event := templateEvent.Prepare()
	eventJSON, err := event.ToJSON()

	if err != nil {
		t.Error("Expecting ToJSON err to be nil", err)
	}

	eventJSONExpected := "{\"name\":\"event:name\",\"version\":\"1\",\"flowId\":\"teste-de-id\",\"id\":\"teste-de-id\",\"payload\":null,\"metadata\":{},\"identity\":{},\"auth\":{}}"

	if eventJSON != eventJSONExpected {
		t.Error("Expecting event json to be compliant with json signature", eventJSON, eventJSONExpected)
	}

	fromJSON, err := FromJSON(eventJSON)
	if err != nil {
		t.Error("Expecting FromJSON err to be nil", err)
	}

	if fromJSON.ID != event.ID {
		t.Error("wrongly imported event", fromJSON)
	}
}

func Test_Event_Import(t *testing.T) {
	event := Event{
		Name:    "event:name",
		Version: "1",
		FlowID:  "teste-de-id",
		ID:      "teste-de-id",
	}

	session, err := ImportEventSession(event)

	if err != nil {
		t.Error("Expecting ImportEventSession err to be nil", err)
	}

	if session.SessionID != "teste-de-id" {
		t.Error("Expecting ID to be imported from event", session.SessionID, "teste-de-id")
	}
}
func Test_Event_Import_From_JSON(t *testing.T) {
	eventJSONBase := "{\"name\":\"event:name\",\"version\":\"1\",\"flowId\":\"teste-de-id\",\"id\":\"teste-de-id\",\"payload\":null,\"metadata\":{},\"identity\":{},\"auth\":{}}"

	session, err := ImportJSONEventSession(eventJSONBase)

	if err != nil {
		t.Error("Expecting ImportJSONEventSession err to be nil", err)
	}

	if session.SessionID != "teste-de-id" {
		t.Error("Expecting ID to be imported from eventJSONBase", session.SessionID, "teste-de-id")
	}
}
