package questions

import (
	"sync"
)

type Request struct {
	Id    string
	Value int
}

type RequestManager struct {
	request *Request
	mtx     *sync.Mutex
}

func NewRequestManager() *RequestManager {
	return &RequestManager{mtx: &sync.Mutex{}}
}

func (r *RequestManager) Set(request *Request) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.request = request
}

func (r *RequestManager) Update(request *Request) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if r.request == nil {
		return
	}
	r.request.Id = request.Id
	r.request.Value = request.Value
}

func (r *RequestManager) Delete() {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.request = nil
}
