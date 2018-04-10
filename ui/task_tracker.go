package ui

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

const progressColor = "blue"

var (
	Term              = !color.NoColor
	progressIndicator = spinner.CharSets[39] // spinning hearth
)

type TaskTracker struct {
	s *spinner.Spinner
}

func NewTaskTracker() *TaskTracker {
	return NewTaskTrackerWithTerm(Term)
}

func NewTaskTrackerWithTerm(term bool) *TaskTracker {
	if term {
		s := spinner.New(progressIndicator, 300*time.Millisecond)
		return &TaskTracker{s}
	}
	return &TaskTracker{}
}

func (t *TaskTracker) Start(msg string) {
	if t.isTerm() {
		t.s.Prefix = msg
		t.s.Color(progressColor)
	} else {
		fmt.Println(msg)
	}
}

func (t *TaskTracker) Step(msg string) {
	if t.isTerm() {
		t.s.Prefix = msg
	} else {
		fmt.Println(msg)
	}
}

func (t *TaskTracker) Success(msg string) {
	t.end(fmt.Sprintf("%s  %s\n", msg, DoneCheck()))
}

func (t *TaskTracker) Failure(msg string) {
	t.end(fmt.Sprintf("%s  %s\n", msg, ErrorCheck()))
}

func (t *TaskTracker) end(msg string) {
	if t.isTerm() {
		t.s.FinalMSG = msg
		t.s.Stop()
	} else {
		fmt.Println(msg)
	}
}

func (t *TaskTracker) isTerm() bool {
	return Term && t.s != nil
}
