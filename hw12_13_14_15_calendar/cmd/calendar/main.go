package main

import (
	_ "github.com/lib/pq"
	"github.com/seth2810/otus_homework/hw12_13_14_15_calendar/cmd/calendar/commands"
)

func main() {
	commands.Execute()
}
