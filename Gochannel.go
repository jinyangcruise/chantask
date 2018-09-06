package chantask

import (
	"sync"
)

type Gochannel struct {
	closeSend    chan struct{}
	closeRecv    chan struct{}
	ret          chan interface{}
	wgSenders    sync.WaitGroup
	wgReceivers  sync.WaitGroup
	routineSega  chan struct{}
	senders      []SenderFunc
	receivers    []ReceiverFunc
	senderArgs   []interface{}
	receiverArgs []interface{}
	started      bool
	ended        bool
}

type SenderFunc func(task *Gochannel, args ...interface{})
type ReceiverFunc func(task *Gochannel, args ...interface{})

// First create a task. The task manage a chan by which you can send/receive any data.
// Then add some senders which can send data to this chan and add some receivers to receive these data from it.
// Finally call task.Start() and senders and receivers will run on goroutines automatically.
// You can determine the buffer size of the chan and how many senders can run at the same time at most.
func GoChannelTask(chanBufferSize, routineNum int) *Gochannel {
	if routineNum <= 0 {
		routineNum = 1
	}
	return &Gochannel{
		closeSend:   make(chan struct{}),
		closeRecv:   make(chan struct{}),
		ret:         make(chan interface{}, chanBufferSize),
		routineSega: make(chan struct{}, routineNum),
	}
}

// Please be careful of the `args`. If you use a param in the SenderFunc and the param is from outside of the
// SenderFunc, you'd better pass the param by the `args ...interface{}` unless the param is stable which means you must
// make sure it will not change until task stops.
func (task *Gochannel) AddSender(senderFunc SenderFunc, args ...interface{}) error {
	if task.ended {
		return ErrTaskStoped
	}
	if task.started {
		return ErrTaskIsRunning
	}
	task.senders = append(task.senders, senderFunc)
	task.senderArgs = append(task.senderArgs, args)
	return nil
}

// Please be careful of the `args`. @see AddSender
func (task *Gochannel) AddReceiver(receiverFunc ReceiverFunc, args ...interface{}) error {
	if task.ended {
		return ErrTaskStoped
	}
	if task.started {
		return ErrTaskIsRunning
	}

	task.receivers = append(task.receivers, receiverFunc)
	task.receiverArgs = append(task.receiverArgs, args)
	return nil
}

func (task *Gochannel) Start() error {
	if task.ended {
		return ErrTaskStoped
	}
	if task.started {
		return ErrTaskIsRunning
	}

	task.started = true

	// sender part
	go func() {
		needBlockRoutine := cap(task.routineSega) < len(task.senders)
		for k, produceFn := range task.senders {
			task.wgSenders.Add(1)
			if task.IsSendStopped() {
				task.wgSenders.Done()
				break
			}
			if needBlockRoutine {
				task.routineSega <- struct{}{} //block here when task.routineSega reaches max buffer size
			}
			go func(fn SenderFunc, args []interface{}) {
				defer func() {
					task.wgSenders.Done()
					if needBlockRoutine {
						<-task.routineSega
					}
				}()
				if task.IsSendStopped() {
					return
				}
				fn(task, args...)
			}(produceFn, task.senderArgs[k].([]interface{}))
		}
		task.wgSenders.Wait()
		close(task.ret) // Close task.ret to make sure Senders not blocked when no data in the chan.
	}()

	// receiver part
	for k, recdFn := range task.receivers {
		task.wgReceivers.Add(1)
		go func(fn ReceiverFunc, args []interface{}) {
			defer func() {
				task.wgReceivers.Done()
			}()
			if task.IsReceiveStopped() {
				return
			}
			fn(task, args...)
		}(recdFn, task.receiverArgs[k].([]interface{}))
	}

	task.wgSenders.Wait()
	task.wgReceivers.Wait()
	task.ended = true
	for range task.ret {
	}
	return nil
}

func (task *Gochannel) IsSendStopped() bool {
	select {
	case <-task.closeSend:
		return true
	default:
		return false
	}
}

func (task *Gochannel) Send(data interface{}) (ok bool) {
	select {
	case <-task.closeSend:
		return false
	default:
		task.ret <- data
		return true
	}
}

// When StopSend is called, all the running Senders will
// stop send data and the waiting Senders will not run.
func (task *Gochannel) StopSend() {
	select {
	case <-task.closeSend:
	default:
		close(task.closeSend)
	}
}

func (task *Gochannel) IsReceiveStopped() bool {
	select {
	case <-task.closeRecv:
		return true
	default:
		return false
	}
}

func (task *Gochannel) Receive() (data interface{}, ok bool) {
	select {
	case <-task.closeRecv:
		return nil, false
	default:
		v, ok := <-task.ret
		if !ok {
			return v, false
		}
		return v, true
	}
}

// When StopReceive is called, all the running Receivers will
// receive a nil data and the waring Receivers will not run.
func (task *Gochannel) StopReceive() {
	select {
	case <-task.closeRecv:
	default:
		close(task.closeRecv)
	}
}
