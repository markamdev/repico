package internal

// MessageType type of internal RePiCo message
type MessageType int

const (
	// Undefined is type of unitialized message
	Undefined MessageType = iota
	// Info is an informative message with some payload (string)
	Info
	// HTTPServerError informs about fatal server error
	HTTPServerError
	// ClosedBySignal informs that interrupt signal has been received
	ClosedBySignal
)

// Message is used for sending information from different routines to main function
type Message struct {
	Type    MessageType
	Content string
}

// MessageBus internal communication channel
var MessageBus chan (Message)

func init() {
	// buffer for 2 messages should be enough to not block any go routine
	MessageBus = make(chan Message, 2)
}
