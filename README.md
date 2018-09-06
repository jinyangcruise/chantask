# chantask

## what is chantask?

chantask is a task ;).

## when to use chantask?

Imagine this: 

You are a police chief and you send all your guys to find Trump's lost mind.
They are sent to different cities and countries all over the world.
All the police should stop finding and return home as long as one find a Trump's lost mind because travel costs much.
If the story is a golang coding task, you must use goroutine to let the police act by themselves, however, 
how to stop goroutines when you have already achieved your task? Let chantask help you.

## how to use chantask?

pseudocode of the story above using chantask should be:

```
task := Create a chantask()

task.AddSender(func () { police0 to find in Trump's bathroom. IF `find one` THEN `task.Send(it)`. })
task.AddSender(func () { police1 to find on Twitter.          IF `find one` THEN `task.Send(it)`. })
task.AddSender(func () { police2 to find in Korea.            IF `find one` THEN `task.Send(it)`. })
task.AddSender(func () { police3 to find in Europe.           IF `find one` THEN `task.Send(it)`. })
task.AddSender(func () { police4 to find in China.            IF `find one` THEN `task.Send(it)`. })
...

task.AddReceiver(func () {  LOOP : IF `the task.Receive() seems not that bad` THEN `push it to Trump's head && task.StopSend() && break` ELSE `continue`. })
task.AddReceiver(func () { IF `runs out of money || Trump's term of service is over` THEN `task.StopSend() && task.StopReceive()`. })

task.Start()
```

