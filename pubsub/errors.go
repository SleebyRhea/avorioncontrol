package pubsub

type ErrUnexpectedMessageType struct {
	subid string
}

func (e *ErrUnexpectedMessageType) Error() string {
	return "Unexpected message type received on: " + e.subid
}

type ErrMessageTimedOut struct {
	subid string
}

func (e *ErrMessageTimedOut) Error() string {
	return "message failed to send due to a timeout on: " + e.subid
}
