package events

type Logger interface {
	Log(log string)
}

type Event interface {
	GetName() string
}

type Handler interface {
	GetEventName() string
	Handle(interface{}) error
}

type Manager struct {
	Handlers []Handler
	Logger   Logger
}

func NewEventManager(Logger Logger) *Manager {
	return &Manager{
		Logger: Logger,
	}
}

func (m *Manager) Dispatch(event Event) {
	for _, handler := range m.Handlers {
		if event.GetName() != handler.GetEventName() {
			continue
		}

		err := handler.Handle(event)
		if err != nil {
			m.Logger.Log(err.Error())
		}
	}
}

func (m *Manager) RegisterHandler(handler Handler) {
	handlers := append(m.Handlers, handler)
	m.Handlers = handlers
}
