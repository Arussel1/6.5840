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

	"rpc.go"
)

type RPCType int
const (
    AskTask RPCType = iota
    ReportDone
)
type RegisterArgs struct {
    Type     RPCType
    WorkerID int

    // for ReportDone
    Task     TaskType
    TaskID   int
}
type RegisterReply struct {
    Task    TaskType
    TaskID  int
    Version int

    // map-only
    File string

    NMap    int
    NReduce int
}

type Coordinator struct {
	mu                sync.Mutex
	mapTasks          []Task
	reduceTasks   	  []Task
	nReduce 	  	  int
	nMap 		      int
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

func (c *Coordinator) AnswerRPC(req *RegisterArgs, res *RegisterReply) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	mapAllDone := true
    reduceAllDone := true
	
	if req.Type != AskTask { return nil }

	for i := 0; i < c.nMap; i++ {
		curMapTask := &c.mapTasks[i]
		if curMapTask.State == Idle {
			curMapTask.Version++
			curMapTask.State = InProgress
			curMapTask.WorkerID = req.WorkerID
			curMapTask.StartTime = time.Now()

			res.Task = TaskMap
			res.TaskID = i 
			res.File = curMapTask.FileName
			res.NReduce = c.nReduce
			res.NMap = c.nMap
			res.Version = curMapTask.Version

			return nil
		}
		if curMapTask.State != Completed { mapAllDone = false }
	}
	if !mapAllDone{
		res.Task = TaskWait
		return nil
	}
	for i := 0; i < c.nReduce; i++ {
			curReduceTask := &c.reduceTasks[i]
			if curReduceTask.State == Idle {
				curReduceTask.Version++
				curReduceTask.State = InProgress
				curReduceTask.WorkerID = req.WorkerID
				curReduceTask.StartTime = time.Now()					

				res.Task = TaskReduce
				res.TaskID = i 
				res.NReduce = c.nReduce
				res.NMap = c.nMap
				res.Version = curReduceTask.Version

				return nil
			}
		if curReduceTask.State != Completed { mapAllDreduceAllDoneone = false }
	}

	if !reduceAllDone {
		res.Task = TaskWait
		return nil
	}	
	// no task available, set res to empty here
	res.Task = TaskExit
	return nil	
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
	}
	// Your code here.

	c.server(sockname)
	return &c
}