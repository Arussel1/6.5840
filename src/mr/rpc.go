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
	Success Status = iota
	Failure 
)
const (
	Map TaskType = iota
	Reduce 
	Wait
)
type Args struct {
	TaskID int
	TaskType TaskType
	FileName string
}

type Reply struct {
	Status Status
	Error string
}