package fxk

import (
	"fmt"
	"net"
	"strconv"
)

type Server struct {
	name      string
	ipVersion string
	ip        string
	port      int
	handle    *Handle
	sess      Cmap
	proty     Cmap
	stratcb   func(*Session) error
	stopcb    func(*Session) error
}

func (s *Server) RstratCb(cb func(*Session) error) {
	s.stratcb = cb
}

func (s *Server) RstopCb(cb func(*Session) error) {
	s.stopcb = cb
}

func (s *Server) Stop() {
	s.sess.Clear()
	s.handle.DestroyWrokerPool()
}

func (s *Server) QuickStart() {
	s.start()
	select {}
}

func (s *Server) GetName() string {
	return s.name
}
func (s *Server) GetIpVersion() string {
	return s.ipVersion
}

func (s *Server) SetName(name string) *Server {
	s.name = name
	return s
}
func (s *Server) SetIpVersion(ipVersion string) *Server {
	s.ipVersion = ipVersion
	return s
}
func (s *Server) SetIp(ip string) *Server {
	s.ip = ip
	return s
}
func (s *Server) SetPort(port int) *Server {
	s.port = port
	return s
}

func (s *Server) GetProty() Cmap {
	return s.proty
}

func (s *Server) SetProty(proty Cmap) *Server {
	s.proty = proty
	return s
}

func (s *Server) start() {

	fmt.Printf("Start [%s] ip is %s port is %d \n", s.name, s.ip, s.port)
	go func() {
		//开启wroker工作池
		s.handle.InitWorkerPool()
		addr, err := net.ResolveTCPAddr(s.ipVersion, fmt.Sprintf("%s:%d", s.ip, s.port))
		if err != nil {
			panic(fmt.Sprintf("Resolve TCP Addr error %s", err.Error()))
		}
		listen, err := net.ListenTCP(s.ipVersion, addr)
		if err != nil {
			panic(fmt.Sprintf("ListenTCP TCP error %s", err.Error()))
		}
		index := 0
		for {
			conn, err := listen.AcceptTCP()
			if err != nil {
				fmt.Println("accept tcp is error", err)
				break
			}
			if s.sess.Count() > 3 {
				conn.Close()
				continue
			}
			sess := NewSession(strconv.Itoa(index), conn, s)
			index++
			go sess.Start()
		}

	}()

}
func (s *Server) GetSessCmap() Cmap {
	return s.sess
}
func (s *Server) SetHandle(handle *Handle) *Server {
	s.handle = handle
	return s
}

func New() *Server {
	return &Server{
		name:      "Fink",
		ipVersion: "tcp",
		ip:        "0.0.0.0",
		port:      9002,
		sess:      NewSyncMap(),
		proty:     nil,
		handle:    nil,
	}
}
