package fxk

import (
	"fmt"
	"net"
)

type Session struct {
	id     string
	conn   *net.TCPConn
	closed bool
	done   chan bool
	rwdata chan []byte
	server *Server
}

func NewSession(id string, conn *net.TCPConn, s *Server) *Session {
	sess := &Session{
		id:     id,
		conn:   conn,
		closed: false,
		done:   make(chan bool),
		server: s,
		rwdata: make(chan []byte, 512),
	}
	s.sess.Set(sess.id, sess)
	return sess
}
func (s *Session) write() {
	fmt.Println("Writer is runing")
	defer fmt.Println(s.conn.RemoteAddr().String(), "[conn Writer exit]")
	for {
		select {
		case data := <-s.rwdata:
			if _, err := s.conn.Write(data); err != nil {
				fmt.Println("send rwdata message is error", err)
				return
			}
		case <-s.done:
			return
		}
	}
}

func (s *Session) read() {
	fmt.Println("Reader is runing")
	defer s.Stop()
	defer fmt.Println(s.conn.RemoteAddr().String(), "[conn Reader exit]")
	for {
		buf := make([]byte, 512)
		cnt, err := s.conn.Read(buf)
		if err != nil {
			fmt.Println("read tcp is error", err)
			break
		}
		tag, data, err := Decode(buf[:cnt])
		if err != nil {
			fmt.Println("decode message error", err)
			continue
		}
		message := NewMessage(tag, data)
		req := NewRequest(s, message)
		s.server.handle.ToTask(req)
		// go s.handle.Schedule(req)
	}
}
func (s *Session) Start() {
	go s.read()
	go s.write()
	if s.server.stopcb != nil {
		err := s.server.stratcb(s)
		if err != nil {
			fmt.Println("call back start is fail", err)
		}
	}
}
func (s *Session) Stop() {
	if !s.closed {
		if s.server.stopcb != nil {
			err := s.server.stopcb(s)
			if err != nil {
				fmt.Println("call back stop is fail", err)
			}
		}
		s.closed = true
		s.server.sess.Remove(s.id)
		s.conn.Close()
		s.done <- true
		close(s.done)
		close(s.rwdata)
	}
}

func (s *Session) GetID() string {
	return s.id
}

func (s *Session) GetConn() *net.TCPConn {
	return s.conn
}

func (s *Session) GetClosed() bool {
	return s.closed
}
func (s *Session) GetDone() chan<- bool {
	return s.done
}
func (s *Session) SetDone(done chan bool) *Session {
	s.done = done
	return s
}
func (s *Session) SetClosed(closed bool) *Session {
	s.closed = closed
	return s
}

func (s *Session) SetConn(conn *net.TCPConn) *Session {
	s.conn = conn
	return s
}

func (s *Session) SetRwData(data chan []byte) *Session {
	s.rwdata = data
	return s
}

func (s *Session) SetServer(server *Server) *Session {
	s.server = server
	return s
}

func (s *Session) SendMessage(tag uint32, data string) error {
	b, err := Encode(tag, data)
	if err != nil {
		return err
	}
	s.rwdata <- b
	return nil
}
