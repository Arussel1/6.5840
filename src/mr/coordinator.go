package mr

import (
	"errors"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
	"time"
)

type Task struct {
	ID        int
	TaskType  TaskType
	FileName  string
	State     Status
	StartTime time.Time
	Version   int
}

const TIMEOUT = 15 * time.Second
const (
	Idle Status = iota
	InProgress
	Completed
)

type Phase int
type IntermediateTaskPointer struct {
	workerAddr string
	fileId int
	attempt int
}
type Worker struct {
	lastHeartbeat int // Time ?
	alive bool
	addr string
}
const (
	PhaseMap Phase = iota
	PhaseReduce
	PhaseFinished
)

type Coordinator struct {
	mu            sync.Mutex
	mapTasks      []Task
	reduceTasks   []Task
	nReduce       int
	mapOut        [][]IntermediateTaskPointer
	workers        map[int]Worker
	timeoutPolicy time.Duration
	currentPhase  Phase 
}

// Your code here -- RPC handlers for the worker to call.

// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

// start a thread that listens for RPCs from worker.go
func (c *Coordinator) server(sockname string) {
	rpc.Register(c)
	rpc.HandleHTTP()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatalf("listen error %s: %v", sockname, e)
	}
	go http.Serve(l, nil)
}

// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
func (c *Coordinator) Done() bool {
	ret := false
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.currentPhase == PhaseFinished {
		ret = true
	}

	return ret
}

func (c *Coordinator) RequestTask(args *RequestTaskArgs, reply *RequestTaskReply) error {
	if args == nil {
		return errors.New("nil args")
	}
	if reply == nil {
		return errors.New("nil reply")
	}

	return nil
}

// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeCoordinator(sockname string, files []string, nReduce int) *Coordinator {
	c := Coordinator{}
	c.currentPhase = PhaseMap
	c.timeoutPolicy = TIMEOUT
	c.nReduce = nReduce
	for i := range files {
		mapTask := &Task{
			ID:        i,
			TaskType:  TaskMap,
			FileName:  files[i],
			State:     Idle,
			StartTime: time.Time{},
			Version:   0,
		}
		c.mapTasks = append(c.mapTasks, *mapTask)

	// Your code here.

	c.server(sockname)
	return &c
}
