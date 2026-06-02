package mr

import "time"

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

const (
	PhaseMap Phase = iota
	PhaseReduce
	PhaseFinished
)


// Add your RPC definitions here.
type Status int
type TaskType int

const (
	TaskMap TaskType = iota
	TaskReduce
	TaskWait
	TaskExit
)

type RequestTaskArgs struct {
	WorkerID int
}

type RequestTaskReply struct {
	TaskID   int
	TaskType TaskType
	FileName string
	NMap     int
	NReduce  int
	Version  int
}

type ReportTaskDoneArgs struct {
	TaskType TaskType
	TaskID   int
	Version  int
	WorkerID int
}

type ReportTaskDoneReply struct {
	Accepted bool
}
