package events

import (
	"encoding/json"
)

func mergeMaps(mergeable ...map[string]interface{}) map[string]interface{} {
	output := make(map[string]interface{})
	for _, m := range mergeable {
		for k, v := range m {
			output[k] = v
		}
	}

	return output
}

// SessionIDGenerator is a basic contract to be
// used as an session identifier generator
type SessionIDGenerator func() string

// EventSession is the base structure for sessions
// to be used in a couple of events and to be
// trace through its sessionID
type EventSession struct {
	Identity   map[string]interface{}
	Metadata   map[string]interface{}
	Auth       map[string]interface{}
	SessionID  string
	Events     map[string]interface{}
	GenerateID SessionIDGenerator
}

// EventTemplate is used to create a base template
// to be easyly replicated
type EventTemplate struct {
	session  *EventSession
	name     string // Preciso de um validador de nome
	version  string
	payload  map[string]interface{}
	metadata map[string]interface{}
	identity map[string]interface{}
	auth     map[string]interface{}
}

// Event is the main package, applying Guiabolso's
// events protocol
type Event struct {
	Name     string                 `json:"name"` // Preciso de um validador de nome
	Version  string                 `json:"version"`
	FlowID   string                 `json:"flowId"`
	ID       string                 `json:"id"`
	Payload  map[string]interface{} `json:"payload"`
	Metadata map[string]interface{} `json:"metadata"`
	Identity map[string]interface{} `json:"identity"`
	Auth     map[string]interface{} `json:"auth"`
}

// ToJSON directly converts event into a JSON package
func (e Event) ToJSON() (string, error) {
	eventJSON, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	eventJSONString := string(eventJSON)
	return eventJSONString, nil
}

// FromJSON directly converts JSON package into an Event
// Use it only if you just want the Event object
// If you want to handle the event, use ImportEventSession instead
func FromJSON(eventJSON string) (Event, error) {
	var event Event
	err := json.Unmarshal([]byte(eventJSON), &event)
	return event, err
}

// GenerateEventSession generates an session
func GenerateEventSession(sessionIDGenerator SessionIDGenerator) EventSession {
	eventSession := &EventSession{
		SessionID:  sessionIDGenerator(),
		GenerateID: sessionIDGenerator,
		Auth:       make(map[string]interface{}),
		Identity:   make(map[string]interface{}),
		Metadata:   make(map[string]interface{}),
	}
	eventSession.Events = make(map[string]interface{})

	return *eventSession
}

// ImportJSONEventSession use an JSON event to generate
// the ID and reuse maximum of its resources as flowID and
// ID. Also registers an event with given JSON event name and version
func ImportJSONEventSession(eventJSON string) (EventSession, error) {
	var event Event
	var session EventSession
	err := json.Unmarshal([]byte(eventJSON), &event)

	if err != nil {
		return session, err
	}

	session, err = ImportEventSession(event)
	return session, err
}

// ImportEventSession use an event to generate
// the ID and reuse maximum of its resources as flowID and
// ID. Also registers an event with given event name and version
func ImportEventSession(event Event) (EventSession, error) {
	var session EventSession

	flowIDGenertor := func() string {
		return event.FlowID
	}

	session = EventSession{
		SessionID:  event.ID,
		GenerateID: flowIDGenertor,
		Auth:       event.Auth,
		Identity:   event.Identity,
		Metadata:   event.Metadata,
	}

	session.Events = make(map[string]interface{})

	registeredEvent := session.RegisterEvent(event.Name, event.Version)
	registeredEvent.WithPayload(event.Payload)

	return session, nil
}

// SetIdentity attachs identity to be used with given session
func (session *EventSession) SetIdentity(identity map[string]interface{}) *EventSession {
	session.Identity = identity
	return session
}

// SetAuth attachs identity to be used with given session
func (session *EventSession) SetAuth(auth map[string]interface{}) *EventSession {
	session.Auth = auth
	return session
}

// SetMetadata is used to attach a metadata into an registered
// Event. This metadata will be reused until overrided
func (session *EventSession) SetMetadata(metadata map[string]interface{}) *EventSession {
	session.Metadata = metadata
	return session
}

// RegisterEvent registers events with given session
// an registered event must be prepare to be puse
func (session *EventSession) RegisterEvent(eventName string, version string) *EventTemplate {
	event := &EventTemplate{
		session: session,
		name:    eventName,
		version: version,
	}

	session.Events[eventName] = event

	return event
}

// WithPayload is used to attach a payload into an registered
// Event. This payload will be reused until overrided
func (e *EventTemplate) WithPayload(payload map[string]interface{}) *EventTemplate {
	e.payload = payload
	return e
}

// Prepare will convert a template event
// into an event object
// also, prepare will tie ID, FlowID and other session
// based attributes into that event
func (e *EventTemplate) Prepare() Event {
	event := Event{
		ID:       e.session.SessionID,
		FlowID:   e.session.GenerateID(),
		Name:     e.name,
		Version:  e.version,
		Payload:  e.payload,
		Identity: mergeMaps(e.session.Identity, e.identity),
		Metadata: mergeMaps(e.session.Metadata, e.metadata),
		Auth:     mergeMaps(e.session.Auth, e.auth),
	}
	return event
}
