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

// Add your RPC definitions here.
type Status int
type TaskType int
const (
	Idle Status = iota
	Success 
	Failure
)
const (
	Map TaskType = iota
	Reduce 
	Wait
	Exit
)
type Task struct {
	ID int
	TaskType TaskType
	State Status
	StartTime time.Time
	FileName string
}

type Reply struct {
	Status Status
	Error string
}