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

func (c *Coordinator) AskTask(args *AskTaskArgs, reply *AskTaskReply) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.resetTimedOutTasks()

    if !c.allMapsDone() {
        c.assignTask(c.mapTasks, TaskMap, reply)
        return nil
    }

    if !c.allReducesDone() {
        c.assignTask(c.reduceTasks, TaskReduce, reply)
        return nil
    }

    reply.TaskType = TaskExit
    return nil
}

func (c *Coordinator) ReportTaskDone(args *ReportTaskDoneArgs, reply *ReportTaskDoneReply) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    var tasks []Task

    if args.TaskType == TaskMap {
        tasks = c.mapTasks
    } else if args.TaskType == TaskReduce {
        tasks = c.reduceTasks
    } else {
        reply.Accepted = false
        return nil
    }

    if args.TaskID < 0 || args.TaskID >= len(tasks) {
        reply.Accepted = false
        return nil
    }

    task := &tasks[args.TaskID]

    if task.State == InProgress && task.Version == args.Version {
        task.State = Completed
        reply.Accepted = true
        return nil
    }

    reply.Accepted = false
    return nil
}

func (c *Coordinator) assignTask(tasks []Task, taskType TaskType, reply *AskTaskReply) {
    for i := range tasks {
        task := &tasks[i]

        if task.State == Idle {
            task.State = InProgress
            task.StartTime = time.Now()
            task.Version++

            reply.TaskType = taskType
            reply.TaskID = task.ID
            reply.Filename = task.FileName
            reply.NMap = c.nMap
            reply.NReduce = c.nReduce
            reply.Version = task.Version
            return
        }
    }

    reply.TaskType = TaskWait
}


func (c *Coordinator) resetTimedOutTasks() {
    now := time.Now()

    for i := range c.mapTasks {
        task := &c.mapTasks[i]
        if task.State == InProgress && now.Sub(task.StartTime) > TIMEOUT {
            task.State = Idle
        }
    }

    for i := range c.reduceTasks {
        task := &c.reduceTasks[i]
        if task.State == InProgress && now.Sub(task.StartTime) > TIMEOUT {
            task.State = Idle
        }
    }
}

func (c *Coordinator) allMapsDone() bool {
	    for i := 0; i < c.nMap; i++ {
        if c.mapTasks[i].State != Completed {
            return false
        }
    }
	return true
}

func (c *Coordinator) allReducesDone() bool {
	    for i := 0; i < c.nReduce; i++ {
        if c.reduceTasks[i].State != Completed {
            return false
        }
    }
	return true
}

// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
func (c *Coordinator) Done() bool {
    c.mu.Lock()
    defer c.mu.Unlock()

    return c.allReducesDone()
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

// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeCoordinator(sockname string, files []string, nReduce int) *Coordinator {

	lenFiles := len(files)

	c := Coordinator{
		nMap: lenFiles,
		nReduce: nReduce,
		mapTasks: make([]Task, lenFiles),
		reduceTasks: make([]Task, nReduce),
	}

	for i := 0; i < lenFiles; i++ {
		c.mapTasks[i] = Task{
			ID: i,
			State: Idle,
			FileName: files[i],
			TaskType: TaskMap,
		}
	}

	for i := 0; i < nReduce; i++ {
		c.reduceTasks[i] = Task {
			ID: i,
			State: Idle,
			TaskType: TaskReduce,
		}
	}

	c.server(sockname)
	return &c
}