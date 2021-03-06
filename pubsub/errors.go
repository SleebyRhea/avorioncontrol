package pubsub

type ErrUnexpectedMessageType struct {
}

func (e *ErrUnexpectedMessageType) Error  {
	
}

type ErrMessageTimedOut struct {
	subid string
}

func (e *ErrMessageTimedOut) Error() string {
	return "message failed to send due to a timeout on: "+e.subid
}