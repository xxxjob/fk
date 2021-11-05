package fxk

type Request struct {
	sess    *Session
	message *Message
}

func (r *Request) GetSess() *Session {
	return r.sess
}
func (r *Request) GetMessage() *Message {
	return r.message
}

func NewRequest(sess *Session, message *Message) *Request {
	return &Request{
		sess:    sess,
		message: message,
	}
}
