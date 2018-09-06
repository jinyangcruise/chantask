package chantask

import "errors"

var ErrTaskIsRunning = errors.New("task is running")
var ErrTaskStoped = errors.New("task stopped")
var ErrTaskNotStart = errors.New("task not start")
