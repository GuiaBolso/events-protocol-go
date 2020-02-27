package events

import "sync"

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
}

// Event is the main package, applying Guiabolso's
// events protocol
type Event struct {
	Name     string // Preciso de um validador de nome
	Version  string
	FlowID   string // Preciso de um gerador
	ID       string // Preciso de um gerador
	Payload  map[string]interface{}
	Metadata map[string]interface{}
	Identity map[string]interface{}
	Auth     map[string]interface{}
}

var (
	once         sync.Once
	eventSession *EventSession
)

// RetrieveEventSession retrieves an existing session or create a new one
func RetrieveEventSession(sessionIDGenerator ...SessionIDGenerator) EventSession {
	once.Do(func() {
		if len(sessionIDGenerator) == 0 {
			panic("Without a sessionIDgenerator, Event is useless")
		}

		generateID := func() string {
			var sessionID string
			for _, c := range sessionIDGenerator {
				sessionID += c()
			}
			return sessionID
		}

		eventSession = &EventSession{
			SessionID:  generateID(),
			GenerateID: generateID,
		}
		eventSession.Events = make(map[string]interface{})
	})

	return *eventSession
}

// WithIdentity attachs identity to be used with given session
func (session *EventSession) WithIdentity(identity map[string]interface{}) *EventSession {
	session.Identity = identity
	return session
}

// WithMetadata is used to attach a metadata into an registered
// Event. This metadata will be reused until overrided
func (session *EventSession) WithMetadata(metadata map[string]interface{}) *EventSession {
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
		Metadata: e.session.Metadata,
	}
	return event
}
