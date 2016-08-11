package task

import (
	"time"
)

type Job interface {
	Execute()
}

type Trigger struct {
	repeatCount int
	repeatInterval time.Duration
	
	startDelay time.Duration
	endDelay time.Duration
}

func (t *Trigger) SetRepeat(repeatCount int, repeatInterval time.Duration) {
	t.repeatCount = repeatCount
	t.repeatInterval = repeatInterval
}

func (t *Trigger) SetDelay(startDelay time.Duration, endDelay time.Duration) {
	t.startDelay = startDelay
	t.endDelay = endDelay
}

type Task struct {
	name string
	trigger *Trigger
	job Job
	
	counter int
	pause bool
	stop bool
	
	next time.Time
}

func NewTask(name string, trigger *Trigger, job Job) *Task {
	t := new(Task)
	
	t.name = name
	t.trigger = trigger
	t.job = job
	
	t.counter = 0
	t.pause = false
	t.stop = true
	
	t.next = time.Now().Local()
	
	return t
}

func (t *Task) GetName() string {
	return t.name
}

func (t *Task) GetJob() Job {
	return t.job
}

func (t *Task) GetTrigger() *Trigger {
	return t.trigger
}

func (t *Task) run() {
	
	counter := 0

	job := t.GetJob()
	trigger := t.GetTrigger()
	
	if trigger.startDelay > 0 {
		time.Sleep(time.Second * trigger.startDelay)
	}
	
	if t.Match() && job != nil {
		job.Execute()
		counter ++
	}
}

func (t *Task) Start() {
	t.pause = false
}


func (t *Task) Pause() {
	t.pause = true
}

func (t *Task) Match() bool {
	now := time.Now().Local()
	
	if t.pause {
		return false
	}
	
	trigger := t.GetTrigger()
	if t.next.Before(now) { 
		if t.counter < trigger.repeatCount {
			t.next = now.Add(time.Second * trigger.repeatInterval)
			t.counter ++;
			return true
		} else {
			t.stop = true
		}
	}
	return false
}
