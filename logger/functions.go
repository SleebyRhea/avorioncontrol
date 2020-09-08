package logger

// sendToChans sends the given string to the provided list of channels
func sendToChans(out string, chs []chan []byte) {
	for _, ch := range chs {
		select {
		case ch <- []byte(out):
		default:
		}
	}
}
