package worker

import (
	"diogonicoleti/strauss/task"
	"fmt"
	"log"
	"time"

	"github.com/Workiva/go-datastructures/queue"
	"github.com/google/uuid"
)

type Worker struct {
	Name      string
	Queue     queue.Queue
	Db        map[uuid.UUID]*task.Task
	TaskCount int
}

func (w *Worker) CollectStats() {
	fmt.Println("I will collect stats")
}

func (w *Worker) RunTask() {
	fmt.Println("I will start or stop a task")
}

func (w *Worker) StartTask(t task.Task) task.DockerResult {
	d := w.newDocker(t)
	t.StartTime = time.Now().UTC()

	result := d.Run()
	if result.Error != nil {
		log.Printf("Error running task: %v: %v\n", t.ID, result.Error)
		t.State = task.Failed
	} else {
		t.ContainerID = result.ContainerId
		t.State = task.Running
	}

	w.Db[t.ID] = &t
	return result
}

func (w *Worker) StopTask(t task.Task) task.DockerResult {
	d := w.newDocker(t)

	result := d.Stop(t.ContainerID)
	if result.Error != nil {
		log.Printf("Error stopping container: %v: %v\n", t.ContainerID,
			result.Error)
	}

	t.FinishTime = time.Now().UTC()
	t.State = task.Completed
	w.Db[t.ID] = &t
	log.Printf("Stopped and removed container %v for task %v\n",
		t.ContainerID, t.ID)

	return result
}

func (w *Worker) newDocker(t task.Task) *task.Docker {
	return task.NewDocker(task.NewConfig(&t))
}
