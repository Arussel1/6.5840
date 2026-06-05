package mr

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net/rpc"
	"os"
	"sort"
	"time"
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

		tempFileName := fmt.Sprintf("temp-%d-%d-*.json", MapIndex, i)
		f, err := os.CreateTemp(".", tempFileName)
		if err != nil {
			log.Printf("Create file: %v", err)
			return false
		}

		encoder := json.NewEncoder(f)
        for _, kv := range buckets[i] {
            err = encoder.Encode(&kv)
            if err != nil {
                f.Close()
                os.Remove(f.Name())
                log.Printf("Encode key/value: %v", err)
                return false
            }
        }

        tempName := f.Name()

        err = f.Close()
        if err != nil {
            os.Remove(tempName)
            log.Printf("Close file: %v", err)
            return false
        }
		

        newName := fmt.Sprintf("mr-%d-%d", MapIndex, i)

        err = os.Rename(tempName, newName)
        if err != nil {
            os.Remove(tempName)
            log.Printf("Rename file: %v", err)
            return false
        }

        log.Printf("File created: %s", newName)
    }
	return true
}

func executeReduce(
    reducefunction func(string, []string) string,
    reduceIndex int,
    split int,
) bool {
    var keyValueList []KeyValue

    for i := 0; i < split; i++ {
        name := fmt.Sprintf("mr-%d-%d", i, reduceIndex)

        file, err := os.Open(name)
        if err != nil {
            log.Printf("Open file %s: %v", name, err)
            return false
        }

        decoder := json.NewDecoder(file)
        for {
            var kv KeyValue
            err := decoder.Decode(&kv)
            if err == io.EOF {
                break
            }
            if err != nil {
                file.Close()
                log.Printf("Decode file %s: %v", name, err)
                return false
            }

            keyValueList = append(keyValueList, kv)
        }

        if err := file.Close(); err != nil {
            log.Printf("Close file %s: %v", name, err)
            return false
        }
    }

    sort.Slice(keyValueList, func(i int, j int) bool {
        return keyValueList[i].Key < keyValueList[j].Key
    })

    tempFileName := fmt.Sprintf("temp-out-%d-*", reduceIndex)

    f, err := os.CreateTemp(".", tempFileName)
    if err != nil {
        log.Printf("Create file %s: %v", tempFileName, err)
        return false
    }

    i := 0
    for i < len(keyValueList) {
        j := i + 1

        for j < len(keyValueList) && keyValueList[j].Key == keyValueList[i].Key {
            j++
        }

        var values []string
        for k := i; k < j; k++ {
            values = append(values, keyValueList[k].Value)
        }

        result := reducefunction(keyValueList[i].Key, values)

        if _, err := fmt.Fprintf(f, "%v %v\n", keyValueList[i].Key, result); err != nil {
            os.Remove(tempFileName)
            log.Printf("Write reduce output: %v", err)
            return false
        }

        i = j
    }

    tempName := f.Name()

    if err := f.Close(); err != nil {
		os.Remove(tempName)
        log.Printf("Close file %s: %v", tempName, err)
        return false
    }

    newName := fmt.Sprintf("mr-out-%d", reduceIndex)

    if err := os.Rename(tempName, newName); err != nil {
		os.Remove(tempName)
        log.Printf("Rename %s to %s: %v", tempName, newName, err)
        return false
    }

    return true
}



func Worker(
    sockname string,
    mapf func(string, string) []KeyValue,
    reducef func(string, []string) string,
) {
    coordSockName = sockname

    for {
        askArgs := AskTaskArgs{}
        askReply := AskTaskReply{}

        if !call("Coordinator.AskTask", &askArgs, &askReply) {
            log.Printf("Worker exit")
            return
        }
        switch askReply.TaskType {
        case TaskWait:
            time.Sleep(time.Second)

        case TaskExit:
            return

        case TaskMap:
            ok := executeMap(mapf, askReply.Filename, askReply.NReduce, askReply.TaskID)
            if ok {
				
                reportDoneArgs := ReportTaskDoneArgs{
                    TaskType: TaskMap,
                    TaskID:   askReply.TaskID,
                    Version:  askReply.Version,
                }
                reportDoneReply := ReportTaskDoneReply{}

                call("Coordinator.ReportTaskDone", &reportDoneArgs, &reportDoneReply)
            }

        case TaskReduce:
            ok := executeReduce(reducef, askReply.TaskID, askReply.NMap)
            if ok {
                reportDoneArgs := ReportTaskDoneArgs{
                    TaskType: TaskReduce,
                    TaskID:   askReply.TaskID,
                    Version:  askReply.Version,
                }
                reportDoneReply := ReportTaskDoneReply{}

                call("Coordinator.ReportTaskDone", &reportDoneArgs, &reportDoneReply)
            }
        }
    }
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
