package fxk

type Router interface {
	PreHandle(req *Request)
	Handle(req *Request)
	PostHandle(req *Request)
}

type AbstrcutRouter struct{}

func (r *AbstrcutRouter) PreHandle(req *Request)  {}
func (r *AbstrcutRouter) Handle(req *Request)     {}
func (r *AbstrcutRouter) PostHandle(req *Request) {}
