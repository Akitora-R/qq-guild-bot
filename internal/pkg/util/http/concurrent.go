package http

import (
	"sort"
	"sync"
)

type BatchResult struct {
	Url      string
	Data     []byte
	Err      error
	RespCode int
	Order    int
}

func BatchExec(requests []Request) []BatchResult {
	wg := &sync.WaitGroup{}
	appendLock := sync.Mutex{}
	var results []BatchResult
	for i, request := range requests {
		wg.Add(1)
		go func(r Request, order int, wg *sync.WaitGroup) {
			bytes, rs, err := r.Exec()
			execResult := BatchResult{
				Url:   r.Url,
				Data:  bytes,
				Err:   err,
				Order: order,
			}
			if rs != nil {
				execResult.RespCode = rs.StatusCode
			}
			appendLock.Lock()
			results = append(results, execResult)
			defer func() {
				wg.Done()
				appendLock.Unlock()
			}()
		}(request, i, wg)
	}
	wg.Wait()
	sort.Slice(results, func(i, j int) bool {
		return results[i].Order < results[j].Order
	})
	return results
}

func BatchExecWithEachCallback(requests []Request, callback func(BatchResult)) {
	for i, request := range requests {
		go func(r Request, order int) {
			bytes, rs, err := r.Exec()
			execResult := BatchResult{
				Data:     bytes,
				Err:      err,
				RespCode: rs.StatusCode,
				Order:    order,
			}
			callback(execResult)
		}(request, i)
	}
}
