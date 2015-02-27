package schedule

import (
	"sync"
	"time"
)

	"github.com/intelsdilabs/pulse/core"
)

type Task struct {
	schResponseChan chan ScheduleResponse
	killChan        chan struct{}
	schedule        Schedule
	workflow    Workflow
	metricTypes []core.MetricType
	mu              sync.Mutex //protects state
	state           TaskState
	creationTime    time.Time
	lastFireTime    time.Time
}

type TaskState int

const (
	//Task states
	TaskStopped TaskState = iota
	TaskSpinning
	TaskFiring
)

func NewTask(s Schedule) *Task {
	return &Task{
		schResponseChan: make(chan ScheduleResponse),
		killChan:        make(chan struct{}),
		schedule:        s,
		state:           TaskStopped,
		creationTime:    time.Now(),
		lastFireTime:    time.Now(),
		workflow: NewWorkflow(),
	}
}

func (t *Task) MetricTypes() []core.MetricType {
	return t.metricTypes
}

func (t *Task) Spin() {
	// We need to lock long enough to change state
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.state == TaskStopped {
		t.state = TaskSpinning
		// spin in a goroutine
		go t.spin()
	}
}

func (t *Task) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.state != TaskStopped {
		t.killChan <- struct{}{}
	}
}

func (t *Task) spin() {
	for {
		// Start go routine to wait on schedule
		go t.waitForSchedule()
		// wait here on
		//  schResponseChan - response from schedule
		//  killChan - signals task needs to be stopped
		select {
		case sr := <-t.schResponseChan:
			// If response show this schedule is stil active we fire
			if sr.State() == ScheduleActive {
				t.lastFireTime = time.Now()
				t.fire()
			}
			// TODO stop task on schedule error state or end state
		case <-t.killChan:
			// Only here can it truly be stopped
			t.state = TaskStopped
			break
		}
	}
}

func (t *Task) fire() {
	t.mu.Lock()
	defer t.mu.Unlock()

	// eventual worker manager goes here
	// simulate work by just a 500ms sleep
	t.state = TaskFiring
	//time.Sleep(time.Millisecond * 500)
	//t.state = TaskSpinning

	//routine fires to get work done (and waits and then updates state)
	go func() {
		t.workflow.Start(t, WorkDispatcher)
		t.mu.Lock()
		if t.state == TaskFiring {
			t.state = TaskSpinning
		}
		t.mu.Unlock()
	}()
}

func (t *Task) waitForSchedule() {
	t.schResponseChan <- t.schedule.Wait(t.lastFireTime)
}
