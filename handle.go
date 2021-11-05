package fxk

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type Handle struct {
	handle        map[uint32]Router
	wokerPoolSize int
	taskQueue     []chan *Request
}

func (h *Handle) Schedule(req *Request) {
	router, ok := h.handle[req.GetMessage().GetId()]
	if !ok {
		fmt.Println("api message id is not found", req.GetMessage().GetId())
	}
	router.PostHandle(req)
	router.Handle(req)
	router.PostHandle(req)
}
func (h *Handle) AddRouter(id uint32, router Router) {
	if _, ok := h.handle[id]; ok {
		panic("repeat api id " + strconv.Itoa(int(id)))
	}
	h.handle[id] = router
}

func (h *Handle) AddHandel(handle map[uint32]Router) {
	h.handle = handle
}

//初始化工作池
func (h *Handle) InitWorkerPool() {
	for i := 0; i < h.wokerPoolSize; i++ {
		h.taskQueue[i] = make(chan *Request, 1024)
		go h.start(i, h.taskQueue[i])
	}
}

func (h *Handle) start(index int, reqchan <-chan *Request) {
	fmt.Println("Woker Startting ID :", index)
	for req := range reqchan {
		h.Schedule(req)
	}
}

func (h *Handle) DestroyWrokerPool() {
	wg := sync.WaitGroup{}
	wg.Add(h.wokerPoolSize)
	go func(taskQueue []chan *Request) {
		for _, chans := range taskQueue {
			close(chans)
			fmt.Println("close worker ======>")
		}
		wg.Done()
	}(h.taskQueue)
	wg.Wait()
}

//随机分配消息给任务队列处理
func (h *Handle) ToTask(req *Request) {
	rand.Seed(time.Now().UnixNano())
	h.taskQueue[rand.Intn(h.wokerPoolSize)] <- req
}

func NewHandle() *Handle {
	return &Handle{
		handle:        make(map[uint32]Router),
		wokerPoolSize: runtime.NumCPU() * 2,
		taskQueue:     make([]chan *Request, runtime.NumCPU()*2),
	}
}
