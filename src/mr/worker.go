package mr

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"net/rpc"
	"os"
)

// Map functions return a slice of KeyValue.
type KeyValue struct {
	Key   string
	Value string
}

// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

var coordSockName string // socket for coordinator

func executeMap(
		mapfunction func( string, string) []KeyValue, 
		filename string, 
		nReduce int, 
		MapIndex int, 
) bool {
	
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Read file %s: %v", filename, err)
		return false
	}

	keyValueList := mapfunction(filename, string(content))
	buckets := make([][]KeyValue, nReduce)

	for  _, kv := range keyValueList {
		r := ihash(kv.Key) % nReduce
		buckets[r] = append(buckets[r],kv)
	}
	for i := 0; i < nReduce; i++ {

		tempFileName := fmt.Sprintf("temp-%d-%d*.json", MapIndex, i)
		f, err := os.CreateTemp(".", tempFileName)
		if err != nil {
			log.Printf("Create file: %v", err)
			return false
		}

		encoder := json.NewEncoder(f)
		err = encoder.Encode(buckets[i])
		if err != nil {
			f.Close()
			os.Remove(f.Name())
			log.Printf("Encode file: %v", err)
			return false
		}

		err = f.Close()
		if err != nil {
			os.Remove(f.Name())
			log.Printf("Close file: %v", err)
			return false
		}
		newName := fmt.Sprintf("mr-%d-%d",MapIndex,i)
		err = os.Rename(f.Name(), newName)
		if err != nil {
			os.Remove(f.Name())
			log.Printf("Rename file: %v", err)
			return false
		}
		
		log.Printf("File created: %s", newName)
	}
	return true
}


// main/mrworker.go calls this function.
func Worker(sockname string, mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	coordSockName = sockname

	// Your worker implementation here.

	// uncomment to send the Example RPC to the coordinator.
	// CallExample()

}

// example function to show how to make an RPC call to the coordinator.
//
// the RPC argument and reply types are defined in rpc.go.
func CallExample() {

	// declare an argument structure.
	args := ExampleArgs{}

	// fill in the argument(s).
	args.X = 99

	// declare a reply structure.
	reply := ExampleReply{}

	// send the RPC request, wait for the reply.
	// the "Coordinator.Example" tells the
	// receiving server that we'd like to call
	// the Example() method of struct Coordinator.
	ok := call("Coordinator.Example", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Printf("reply.Y %v\n", reply.Y)
	} else {
		fmt.Printf("call failed!\n")
	}
}

// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	c, err := rpc.DialHTTP("unix", coordSockName)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	if err := c.Call(rpcname, args, reply); err == nil {
		return true
	}
	log.Printf("%d: call failed err %v", os.Getpid(), err)
	return false
}
