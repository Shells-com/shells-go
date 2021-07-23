package main

// Fix for windows builds?

// #cgo windows LDFLAGS: -Wl,-Bstatic -lssp -Wl,-Bdynamic
import "C"
