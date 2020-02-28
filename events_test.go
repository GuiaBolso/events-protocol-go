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

	session := RetrieveEventSession(basicGenerator)

	if session.SessionID != "1" {
		t.Error("Expecting sessionId to be 1:", session.SessionID)
	}

	session2 := RetrieveEventSession()

	if session2.SessionID != session.SessionID {
		t.Error("Expecting session to be a application singleton:", session2.SessionID, session.SessionID)
	}
	count = 0
}

func Test_Event(t *testing.T) {
	eventName := "event:name"
	mockIDGenerator := func() string {
		return fmt.Sprintf("%d", rand.Int())
	}
	session := RetrieveEventSession(mockIDGenerator)

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
