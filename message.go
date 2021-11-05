package fxk

type Message struct {
	id   uint32
	data string
}

func (m *Message) GetData() string {
	return m.data
}
func (m *Message) GetId() uint32 {
	return m.id
}

func NewMessage(id uint32, data string) *Message {
	return &Message{
		id:   id,
		data: data,
	}
}
