package mr

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
