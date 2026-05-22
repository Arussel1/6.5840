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
	InProgress
	Completed
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
	FileName string
	State Status
	StartTime time.Time
	
}

type Reply struct {
	ID int
	TaskType TaskType
	FileName string
	NMap int
	NReduce int
}