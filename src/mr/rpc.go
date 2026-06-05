package mr

import (
	"sync"
	"time"
)

//
// RPC definitions.
//
// remember to capitalize all names.
//

//
// example to show how to declare the arguments
// and reply for an RPC.
//

type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}

type Status int

const (
	Idle Status = iota
	InProgress
	Completed
)

type TaskType int

const (
	TaskMap TaskType = iota
	TaskReduce
	TaskWait
	TaskExit
)


type Task struct {
	mu 	sync.Mutex
	ID        int
	TaskType  TaskType
	FileName  string
	State     Status
	StartTime time.Time
	Version   int
	WorkerID int
}

const TIMEOUT = 15 * time.Second


// Add your RPC definitions here.




type AskTaskArgs struct {}

type AskTaskReply struct {
	TaskType TaskType
	TaskID int
	Filename string
	TaskIndex int 
	NReduce int 
	NMap int
	Version int
}


type ReportTaskDoneArgs struct {
	TaskType TaskType
	TaskID   int
	Version  int
}

type ReportTaskDoneReply struct {
	Accepted bool
}
