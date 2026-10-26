package coroutine

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"sort"
	"sync"
	"time"
)

// simpleTask 简易的协程池(也不算是个完整的池,他只是限制了协程数量)
type simpleTask struct {
	// 上下文
	ctx ctx.Context
	// 任务池大小
	poolSize int
	// 任务
	tasks []func(params ...interface{}) interface{}
	// 任务参数
	taskParams []interface{}
	// 任务执行游标
	taskIndex int
	// 未排序返回值
	result map[int]interface{}
	// 排序的返回值
	sortResult []interface{}
	// 批次执行完是否删除记录
	batchExecutedDeleteRecord bool
	// 批次执行回调
	batchExecuted func(executedTaskLength int, executedResult []interface{}) error
}

func (t *simpleTask) setPoolSize(size int) {
	t.poolSize = size
	t.taskIndex = 0
	t.result = map[int]interface{}{}
	t.batchExecutedDeleteRecord = false
	t.batchExecuted = func(executedTaskLength int, executedResult []interface{}) error {
		return nil
	}
}

func (t *simpleTask) PushTask(task func(params ...interface{}) interface{}, params ...interface{}) {
	t.tasks = append(t.tasks, task)
	t.taskParams = append(t.taskParams, params)
}

func (t *simpleTask) SetBatchExecutedDeleteRecord(isDeleted bool) {
	t.batchExecutedDeleteRecord = isDeleted
}

func (t *simpleTask) SetBatchExecuted(fn func(executedTaskLength int, executedResult []interface{}) error) {
	t.batchExecuted = fn
}

func (t *simpleTask) Run() ([]interface{}, error) {
	var wg sync.WaitGroup
	var lock sync.Mutex
	var err error
	taskLength := len(t.tasks)
	for t.taskIndex < taskLength {
		wgNum := taskLength - t.taskIndex
		if wgNum >= t.poolSize {
			wg.Add(t.poolSize)
		} else {
			wg.Add(wgNum)
		}
		if t.batchExecutedDeleteRecord {
			t.sortResult = []interface{}{}
			t.result = map[int]interface{}{}
		}
		for i := 1; i <= t.poolSize; i++ {
			if t.taskIndex >= taskLength {
				break
			}
			task := t.tasks[t.taskIndex]
			taskParams := t.taskParams[t.taskIndex]
			time.Sleep(time.Second / 10)
			go func(index int, fn func(params ...interface{}) interface{}, p ...interface{}) {
				result := fn(p...)
				lock.Lock()
				t.result[index] = result
				lock.Unlock()
				wg.Done()
			}(t.taskIndex, task, taskParams)
			t.taskIndex = t.taskIndex + 1
		}
		wg.Wait()
		startBatchSortTime := time.Now()
		log.Info(t.ctx).Msgf("开始批次结果排序...%d", t.taskIndex)
		err = t.sort()
		log.Info(t.ctx).Msgf("批次结果排序耗时%s", time.Since(startBatchSortTime))
		if err != nil {
			return nil, err
		}
		err = t.batchExecuted(t.taskIndex, t.sortResult)
		if err != nil {
			return nil, err
		}
	}
	return t.sortResult, nil
}

func (t *simpleTask) sort() error {
	var indexList []int
	var newResult []interface{}
	for key, item := range t.result {
		switch item.(type) {
		case error:
			return item.(error)
		}
		indexList = append(indexList, key)
	}
	sort.Ints(indexList)
	for _, page := range indexList {
		newResult = append(newResult, t.result[page])
	}
	t.sortResult = newResult
	return nil
}

func NewSimpleTask(size int, ctx ctx.Context) *simpleTask {
	task := &simpleTask{}
	task.setPoolSize(size)
	task.ctx = ctx
	return task
}
