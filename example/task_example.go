package main

import (
	"task"
	"fmt"
	"time"
	"sync"
)

type JobA struct {

}

func (j JobA) Execute() {
	fmt.Println("JobA : Now is " + time.Now().String())
}

type JobB struct {

}

func (j JobB) Execute() {
	fmt.Println("JobB : Now is " + time.Now().String())
}

func wait(wg *sync.WaitGroup) chan bool {
	ch := make(chan bool)
	go func() {
		wg.Wait()
		ch <- true
	}()
	return ch
}

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	
	c := task.NewContainer("Test-Container")
	
	t1 := new(task.Trigger)
	t1.SetRepeat(0, 2)
	
	t2 := new(task.Trigger)
	t2.SetRepeat(0, 5)
	t2.SetDelay(5, 0)
	
	c.AddTask("TaskA", t1, new(JobA))
	c.AddTask("TaskB", t2, new(JobB))
	
	fmt.Println("start at " + time.Now().String())
    fmt.Println("----Start (With Daemon)----")
    c.Start()
    
    time.Sleep(time.Second * 11) // pause container
    fmt.Println("---------Stop All Task---------")
    c.StopAllTask()
    
    time.Sleep(time.Second * 10) // restart container
    fmt.Println("--------Resume All Task--------")
    c.ResumeAllTask()
    
    time.Sleep(time.Second * 11) // stop task
    fmt.Println("--------Stop Task A--------")
    c.StopTask("TaskA")
    
    time.Sleep(time.Second * 10) // restart task
    fmt.Println("------Start Task A--------")
    c.ResumeTask("TaskA")
    
    time.Sleep(time.Second * 11) // remove task
    fmt.Println("-----Remove Task A---------")
    c.RemoveTask("TaskA")
    
    time.Sleep(time.Second * 10) // stop all
    fmt.Println("---------Stop Container---------")
    c.StopContainer()
	
	select {
	case <-time.After(time.Second * 10):
		fmt.Println("Finish at " + time.Now().String())
	}

}
