package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sort"
	"os"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// splits a string, trims spaces on every element
func splitandclean(in, sep string, n int) []string {
	v := strings.SplitN(in, sep, n)
	for i := range v {
		v[i] = strings.TrimSpace(v[i])
	}
	return v
}

type Task struct {
	name string
	running bool
	todo map[string]bool
	dependants []*Task
}

var tasks = map[string]*Task{}

func lookup(name string) *Task {
	t := tasks[name]
	if t == nil {
		t = &Task{ name: name, todo: make(map[string]bool) }
		tasks[name] = t
	}
	return t
}

func minready() (*Task, bool) {
	ready := []string{}
	haswaiting := false
	for _, task := range tasks {
		if task.running {
			continue
		}
		haswaiting = true
		if len(task.todo) == 0 {
			ready = append(ready, task.name)
		}
	}
	if len(ready) == 0 {
		return nil, haswaiting
	}
	sort.Strings(ready)
	return tasks[ready[0]], haswaiting
}

func finishtask(task *Task) {
	fmt.Printf("%d %s finished\n", clock, task.name)
	for _, t2 := range task.dependants {
		delete(t2.todo, task.name)
	}
}


var workers []int
var workerjob []*Task
var clock int

func freeworker() int {
	for i := range workerjob {
		if workerjob[i] == nil {
			return i
		}
	}
	return -1
}

func runwork() {
	fmt.Printf("runwork\n")
	mini := -1
	for i := range workers {
		if workerjob[i] == nil {
			continue
		}
		if mini < 0 {
			mini = i
		}
		if workers[i] < workers[mini] {
			mini = i
		}
	}
	n := workers[mini]
	fmt.Printf("working for %d units of time\n", n)
	clock += n
	for i := range workers {
		if workerjob[i] == nil {
			continue
		}
		workers[i] -= n
		if workers[i] == 0 {
			finishtask(workerjob[i])
			workerjob[i] = nil
		}
	}
}

func main() {
	infile := "07.txt"
	isexample := false
	_ = isexample
	if len(os.Args) >= 2 && os.Args[1] == "example" {
		isexample = true
		infile = "07.example"
	}
	buf, err := ioutil.ReadFile(infile)
	must(err)
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := splitandclean(line, " ", -1)
		t1 := lookup(fields[7])
		t1.todo[fields[1]] = true
		t2 := lookup(fields[1])
		t2.dependants = append(t2.dependants, t1)
	}
	
	fmt.Printf("starting\n")
	
	r := []string{}	
	
	var numworkers = 2
	var extratime = 0
	
	if len(tasks) > 7 {
		numworkers = 5
		extratime = 60
	}
	
	workers = make([]int, numworkers)
	workerjob = make([]*Task, numworkers)
	
	mainloop:
	for {
		var i int
		for {
			i = freeworker()
			if i >= 0 {
				break
			}
			runwork()
		}
		var task *Task
		for {
			var haswaiting bool
			task, haswaiting = minready()
			if task != nil {
				break
			}
			if !haswaiting {
				break mainloop
			}
			runwork()
		}
		fmt.Printf("%d starting %s on worker %d\n", clock, task.name, i)
		task.running = true
		r = append(r, task.name)
		workerjob[i] = task
		workers[i] = (int(task.name[0]) - 'A') + 1 + extratime
	}
	
	for {
		busy := false
		for i := range workers {
			if workerjob[i] != nil {
				busy = true
				break
			}
		}
		if !busy {
			break
		}
		runwork()
	}
	part1out := strings.Join(r, "")
	fmt.Printf("PART1: %s\n", part1out)
	fmt.Printf("clock %d\n", clock)
	
	if isexample {
		if part1out != "CAFBDE" || clock != 15 {
			fmt.Printf("FAIL!\n")
		}
	} else {
		if part1out != "GJMFNWVDHBCIETUYQALSPXZORK" || clock != 1050 {
			fmt.Printf("FAIL!\n")
		}
	}
}
