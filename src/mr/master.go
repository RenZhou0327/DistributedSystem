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

const (
	TaskStatusReady   = 0
	TaskStatusQueue   = 1
	TaskStatusRunning = 2
	TaskStatusFinish  = 3
	TaskStatusErr     = 4
)

const (
	MaxTaskRunTime   = time.Second * 5
	ScheduleInterval = time.Millisecond * 500
)

type TaskStat struct {
	Status    int
	WorkerId  int
	StartTime time.Time
}

type Master struct {
	// Your definitions here.
	files     []string
	nReduce   int
	taskPhase TaskPhase
	TaskStats []TaskStat
	mu        sync.Mutex
	done      bool
	workerSeq int
	taskCh    chan Task
}

func (m *Master) getTask(taskSeq int) Task {
	task := Task{
		FileName: "",
		NReduce:  m.nReduce,
		NMaps:    len(m.files),
		Seq:      taskSeq,
		Phase:    m.taskPhase,
		Alive:    true,
	}
	DPrintf("m:%+v, taskseq:%d, lenfiles:%d, lents:%d", m, taskSeq, len(m.files), len(m.taskStats))
	if task.Phase == MapPhase {
		task.FileName = m.files[taskSeq]
	}
	return task
}

// Your code here -- RPC handlers for the worker to call.

//
// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
//
func (m *Master) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

//
// start a thread that listens for RPCs from worker.go
//
func (m *Master) server() {
	rpc.Register(m)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := masterSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrmaster.go calls Done() periodically to find out
// if the entire job has finished.
//
func (m *Master) Done() bool {
	ret := false

	// Your code here.

	return ret
}

//
// create a Master.
// main/mrmaster.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeMaster(files []string, nReduce int) *Master {
	m := Master{}
	m.mu = sync.Mutex{}
	m.nReduce = nReduce
	m.files = files
	if nReduce > len(files) {
		m.taskCh = make(chan Task, nReduce)
	}
	else {
		m.taskCh = make(chan Task, len(m.files))
	}

	// Your code here.

	m.server()
	return &m
}
