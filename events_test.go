package events

import "testing"

func Test_Session(t *testing.T) {
	var basicGenerator SessionIDGenerator
	basicGenerator = func() string {
		return "teste"
	}
	session := RetrieveEventSession(basicGenerator)

	if session.SessionID != "teste" {
		t.Error("Expecting 'test' as sessionId")
	}

	session2 := RetrieveEventSession()

	if session2.SessionID != session.SessionID {
		t.Error("Expecting session to be singleton")
	}
}

func Test_Event(t *testing.T) {
	identity := map[string]interface{}{
		"identity": "identity",
	}
	session := RetrieveEventSession(func() string { return "session-id" })
	session.WithIdentity(identity)

	metadata := map[string]interface{}{
		"metadata": "metadata",
	}
	session.WithMetadata(metadata)

	event := session.RegisterEvent("event:name", "1")

	if len(session.Events) != 1 {
		t.Error("Expecting event list to be 1", session)
	}

	payload := map[string]interface{}{
		"payload": "payload",
	}
	concreteEvent := event.WithPayload(payload).Prepare()

	if concreteEvent.Identity["identity"] != identity["identity"] {
		t.Error("Expecting template identity to be used in final event", concreteEvent)
	}
	if concreteEvent.Metadata["metadata"] != metadata["metadata"] {
		t.Error("Expecting template metadata to be used in final event", concreteEvent)
	}
	if concreteEvent.Payload["payload"] != payload["payload"] {
		t.Error("Expecting template payload to be used in final event", concreteEvent)
	}
}