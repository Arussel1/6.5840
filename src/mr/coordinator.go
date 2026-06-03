package mr

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
	"time"
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
	err := rpc.Register(c)
	if err != nil {
		log.Fatalf("Cannot register server %v", err)
	}
	rpc.HandleHTTP()
	os.Remove(sockname)
	l, err := net.Listen("unix", sockname)
	if err != nil {
		log.Fatalf("listen error %s: %v", sockname, err)
	}
	go func() {
		err := http.Serve(l, nil)
		if err != nil {
		log.Printf("Cannot serve http %v", err)
		}
	}()
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
		if curReduceTask.State != Completed { reduceAllDone = false }
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
	c.mu.Lock()
	defer c.mu.Unlock()
	for i := 0; i < c.nReduce; i++ {
		curReduceTask := &c.reduceTasks[i]
		if curReduceTask.State != Completed {
			return false
		}
	}
	return true
}

// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeCoordinator(sockname string, files []string, nReduce int) *Coordinator {

	lenFiles := len(files)

	c := Coordinator{
		nMap: lenFiles,
		nReduce: nReduce,
		mapTasks: make([]Task,lenFiles),
		reduceTasks: make([]Task,nReduce),
	}

	for i := 0; i < lenFiles; i++ {
		c.mapTasks[i] = Task{
			State: Idle,
			FileName: files[i],
			TaskType: TaskMap,
		}
	}

	for i := 0; i < nReduce; i++ {
		c.reduceTasks[i] = Task {
			State: Idle,
			TaskType: TaskReduce,
		}
	}

	c.server(sockname)
	return &c
}