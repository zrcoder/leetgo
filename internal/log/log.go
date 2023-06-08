package log

// prod mod do nothing
var Trace = func(x ...any) {}
var Tracef = func(format string, x ...any) {}
