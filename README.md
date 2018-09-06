# chantask

make your goroutines run more organized.

## What is chantask?

chantask is a task ;).

## When to use chantask?

Imagine this: 

You are a police chief and you send all your guys to find Trump's lost mind.
They are sent to different cities and countries all over the world.
All the police should stop finding and return home as long as one find a Trump's lost mind because travel costs much.
If the story is a golang coding task, you must use goroutine to let the police act by themselves, however, 
how to stop goroutines when you have already achieved your task? Let chantask help you.

## How to use chantask?

pseudocode of the story above using chantask should be:

```
task := Create a chantask()

task.AddSender(func () { police0 to find in Trump's bathroom. IF `find one` THEN `task.Send(it)`. })
task.AddSender(func () { police1 to find on Twitter.          IF `find one` THEN `task.Send(it)`. })
task.AddSender(func () { police2 to find in Korea.            IF `find one` THEN `task.Send(it)`. })
task.AddSender(func () { police3 to find in Europe.           IF `find one` THEN `task.Send(it)`. })
task.AddSender(func () { police4 to find in China.            IF `find one` THEN `task.Send(it)`. })
...

task.AddReceiver(func () {  
	for {
		IF `mind := task.Receive() seems not that bad` {
			`push the mind to Trump's head`
			`task.StopSend()`
			`break` 
		} 
	}
})
task.AddReceiver(func () { 
	IF `running out of money` || `Trump's term of service is over` {
		`task.StopSend()`
		`task.StopReceive()`.
	}
})

task.Start()
```

## Explanation

Senders are a kind of func that can send data to the chan built in the task.
Receivers are a kind of func that can receive data from the same chan.
Never send datas within the receivers or panic like `data sent to closed chan` will occur.
All the senders and receivers will run by goroutine.

## Example

```
package main

import (
	"github.com/jinyangcruise/chantask"
	"fmt"
)

func main() {
	var f1 chantask.SenderFunc = func(task *chantask.ChanTask, args ...interface{}) {
		id := args[0].(string)
		task.Send(id)
	}

	f2 := func(task *chantask.ChanTask, args ...interface{}) {
		id := args[0].(string)
		task.Send(id)
	}

	var f3 chantask.SenderFunc = func(task *chantask.ChanTask, args ...interface{}) {
		v, ok := task.Receive()
		if ok {
			fmt.Println("f3 got from", v)
		} else {
			fmt.Println("f3 got noting")
		}
	}

	// At most 10 (it's up to you) senders can send data to chan at the same time.
	// At most 10 (it's up to you) senders will be running by goroutine at the same time.
	// If there are over 10 senders, other senders will wait until former sender func return.
	task := chantask.CreateChanTask(10, 10) 

	task.AddSender(f1, "id:1")
	task.AddSender(chantask.SenderFunc(f2), "id:2")

	task.AddReceiver(chantask.ReceiverFunc(f3))
	task.AddReceiver(chantask.ReceiverFunc(func(task *chantask.ChanTask, args ...interface{}) {
		v, ok := task.Receive()
		if ok {
			fmt.Println("f4 got from", v)
		} else {
			fmt.Println("f4 got nothing")
		}
	}))
	task.AddReceiver(chantask.ReceiverFunc(func(task *chantask.ChanTask, args ...interface{}) {
		x := args[0].(int)
		y := args[1].(int)
		v, ok := task.Receive()
		if ok {
			fmt.Println("f5 got from", v, ". sum:", x+y)
		} else {
			fmt.Println("f5 got nothing")
		}

	}), 7, 8)
	task.Start()
}
```
output:
```
f3 got noting
f5 got from id:2 . sum: 15
f4 got from id:1
```
output may be different from above because receivers are competitors for the data in the task.