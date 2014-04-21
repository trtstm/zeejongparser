package throttler

import "net/http"
import "sync"

type workerInput struct {
	url string
	resultChannel chan workerOutput
}

type workerOutput struct {
	resp *http.Response
	err error
}

type ConnectionThrottler struct {
	max int
	nConnections int
	paramChannel chan workerInput
}



func NewConnectionThrottler(max int) *ConnectionThrottler {
	c := ConnectionThrottler{max: max, nConnections: 0}
	c.paramChannel = make(chan workerInput)

	go worker(c.paramChannel, max)

	return &c
}

func (this *ConnectionThrottler) Get(url string) (*http.Response, error) {
	param := workerInput{url: url, resultChannel: make(chan workerOutput)}

	this.paramChannel <- param
	
	result := <- param.resultChannel
	return result.resp, result.err
}

func worker(paramChannel chan workerInput, max int) {
	tasks := []workerInput{}

	executingTasksLock := sync.RWMutex{}
	executingTasks := 0	


	for {
		select {
			case param := <- paramChannel:
				tasks = append(tasks, param)
			default:
				executingTasksLock.Lock()
				if len(tasks) > 0 && executingTasks < max {
					executingTasks += 1
					executingTasksLock.Unlock()

					currentTask := tasks[0]
					tasks = tasks[1:]

					go func(task workerInput) {
						result := workerOutput{}
						result.resp, result.err = http.Get(task.url)
						task.resultChannel <- result

						executingTasksLock.Lock()
						executingTasks -= 1
						executingTasksLock.Unlock()
					}(currentTask)
				} else {
					executingTasksLock.Unlock()
				}
		}
	}
}
