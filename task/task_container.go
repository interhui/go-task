package task

import (
	"time"
	"log"
)

type TaskContainer struct {
	name string
	running bool
	tasks []*Task
	
	add chan *Task
	stop chan bool
}

func NewContainer(name string) *TaskContainer {
	c := new(TaskContainer)
	c.name = name
	c.add = make(chan *Task)
	c.stop = make(chan bool)
	return c
}

func (c *TaskContainer) GetName() string {
	return c.name
}

func (c *TaskContainer) AddTask(taskName string, trigger *Trigger, job Job) bool {
	for _, task := range(c.tasks) {
		if task.GetName() == taskName {
			return false
		}
	}

	task := NewTask(taskName, trigger, job)
	c.tasks = append(c.tasks, task)
	if c.add != nil && c.running == true {
		c.add <- task
	}

	return true
}

func (c *TaskContainer) RemoveTask(taskName string) bool {

	for i := 0; i < len(c.tasks); {
		task := c.tasks[i]
		if task.name ==  taskName {
			task.Stop()
			c.tasks = append(c.tasks[:i], c.tasks[i+1:]...)
			return true
		}
	}

	return false
}

func (c *TaskContainer) GetTask(taskName string) *Task {
	for _, task := range(c.tasks) {
		if task.GetName() == taskName {
			return task
		}
	}
	return nil
}

func (c *TaskContainer) GetTasks() []*Task {
	return c.tasks
}

func (c *TaskContainer) Start() {

	log.Println("Start Container : " + c.GetName())
	
	c.running = true

	go func() {
		for {
			select {
			case <-time.After(time.Second * 1):
				for _, task := range (c.tasks) {
					go task.run()
				}
				continue
			case <-c.stop:
				log.Println("Stop Container : " + c.GetName())
				return
			case task := <-c.add:
				log.Println("Add New Task : " + task.GetName())
				go task.run()
			}
		}
	}()
	
}

func (c *TaskContainer) ResumeAllTask() {
	for _, task := range(c.tasks) {
		task.Start()
	}
}

func (c *TaskContainer) ResumeTask(taskName string) {
	for _, task := range(c.tasks) {
		if taskName == task.name {
			task.Start()
			break
		}
	}
}

func (c *TaskContainer) StopTask(taskName string) {
	for _, task := range(c.tasks) {
		if taskName == task.name {
			task.Stop()
			break
		}
	}
}

func (c *TaskContainer) StopAllTask() {
	for _, task := range(c.tasks) {
		task.Stop()
	}
}

func (c *TaskContainer) StopContainer() {
	c.running = false
	c.stop <- true
}
