Lab 1 MIT distributed system: Distributed Map-Reduce

It is so hard that my whole 2 weeks is stuck on this

First, I tried to blindly jump in, but it didn't work at all, 

After that, I am trying to write a doc of each file on what they did with details functionality

Also need to read the MapReduce paper again

And spent 3 days coding the solution,

first, my map output is all mr-0-%i, which is due to passing TaskIndex vs TaskID

Fix it, and the test worked

Finally with lab1, a lot more to learn ahead

I use AI as a judge to help me unstuck with hint, with caveman prompt "Code Block", eval this, hint code, no code, 
since this lab is usually done in group of 2, so I think it is fair to use AI simulated as a knowledge partner, but still need to be careful of result tho. So far, it is very good.
![Pass test](image.png)


So, I will explain what does each file do (or try my best to)


rpc.go (src/mr/rpc.go)

Store all type of task and Ask/Request including:

Task, 
TIMEOUT const, 
AskTaskArgs/Reply, 
RequestTaskDoneArg/Reply

Details is in the file

coordinator.go (src/mr/coordinator.go)

Coordinator type/construct

AskTask

ReportTaskDone

assignTask

allMapsDone

allReducesDone

resetTimedOutTasks

Done

Will probably explain more when I have tmr, this is the template for now, but yeah, details is in the file

worker.go (src/mr/coordinator.go)

executeMap/ executeReduce: 

Worker constructor

Will probably explain more when I have tmr, this is the template for now, but yeah, details is in the file

Skill I learn:

concurrency handling: if there is a data that could get multiple access, lock it, when in doubtd data could be access concurrently, lock it, race condition is so hard to debug, because it may/ it may not happen

plan -> act -> test -> repeat: For more complex project, jumping in head-on won't work anymore, make sure to plan as details as you can, plan out everything as possible before even writing a single line of code, or else, you will fail, unless you are a genius

Don't repeat Type: Ex: TaskID/TaskIndex, it cause the map to pass as mr-0-X instead of mr-X-X, and it cost me 1 hours to debug

Log + error handling: when in doubt, log, log + debug is so powerful, and help you to catch bug fast, also, err handling first, it is a pattern in go to handle err first, and i understand why this is so helpful now, i will be lost in sea of code if not for this

Keep state simple: My final state is Idle -> In Progress -> Done, which is really simple compared from what I imagined from the start, and the main point is it work, so thing can turn out way simpler that it sounds

ETC.....

