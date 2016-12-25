// gofsm is a simple, featured FSM implementation that has some different features with other FSM implementation.
// One feature of gofsm is it doesn't persist/keep states of objects. When it processes transitions, you must pass current states to id, so you can look gofsm as a "stateless" state machine. This benefit is one gofsm instance can be used to handle transitions of a lot of object instances, instead of creating a lot of FSM instances. Object instances maintain their states themselves.
// Another feature is it provides a common interface for Moore and Mealy FSM. You can implement corresponding methods (OnExit, Action, OnEnter) for those two FSM.
// The third interesting feature is you can export configured transitions into a state diagram. A picture is worth a thousand words.

// Style of gofsm refers to implementation of https://github.com/elimisteve/fsm.

package fsm

import (
	"fmt"
	"os/exec"
	"strings"
)

// Transition is a state transition and all data are literal values that simplifies FSM usage and make it generic.
type Transition struct {
	From   string
	Event  string
	To     string
	Action string
}

// Delegate is used to process actions. Because gofsm uses literal values as event, state and action, you need to handle them with corresponding functions. DefaultDelegate is the default delegate implementation that splits the processing into three actions: OnExit Action, Action and OnEnter Action. you can implement different delegates.
type Delegate interface {
	// HandleEvent handles transitions
	HandleEvent(action string, fromState string, toState string, args []interface{})
}

// StateMachine is a FSM that can handle transitions of a lot of objects. delegate and transitions are configured before use them.
type StateMachine struct {
	delegate    Delegate
	transitions []Transition
}

// Error is an error when processing event and state changing.
type Error interface {
	error
	BadEvent() string
	CurrentState() string
}

type smError struct {
	badEvent     string
	currentState string
}

func (e smError) Error() string {
	return fmt.Sprintf("state machine error: cannot find transition for event [%s] when in state [%s]\n", e.badEvent, e.currentState)
}

func (e smError) BadEvent() string {
	return e.badEvent
}

func (e smError) CurrentState() string {
	return e.currentState
}

// NewStateMachine creates a new state machine.
func NewStateMachine(delegate Delegate, transitions ...Transition) *StateMachine {
	return &StateMachine{delegate: delegate, transitions: transitions}
}

// Trigger fires a event. You must pass current state of the processing object, other info about this object can be passed with args.
func (m *StateMachine) Trigger(currentState string, event string, args ...interface{}) Error {
	trans := m.findTransMatching(currentState, event)
	if trans == nil {
		return smError{event, currentState}
	}

	if trans.Action != "" {
		m.delegate.HandleEvent(trans.Action, currentState, trans.To, args)
	}
	return nil
}

// findTransMatching gets corresponding transition according to current state and event.
func (m *StateMachine) findTransMatching(fromState string, event string) *Transition {
	for _, v := range m.transitions {
		if v.From == fromState && v.Event == event {
			return &v
		}
	}
	return nil
}

// Export exports the state diagram into a file.
func (m *StateMachine) Export(outfile string) error {
	return m.ExportWithDetails(outfile, "png", "dot", "72", "-Gsize=10,5 -Gdpi=200")
}

// ExportWithDetails  exports the state diagram with more graphviz options.
func (m *StateMachine) ExportWithDetails(outfile string, format string, layout string, scale string, more string) error {
	dot := `digraph StateMachine {

	rankdir=LR
	node[width=1 fixedsize=true shape=circle style=filled fillcolor="darkorchid1" ]
	
	`

	for _, t := range m.transitions {
		link := fmt.Sprintf(`%s -> %s [label="%s | %s"]`, t.From, t.To, t.Event, t.Action)
		dot = dot + "\r\n" + link
	}

	dot = dot + "\r\n}"
	cmd := fmt.Sprintf("dot -o%s -T%s -K%s -s%s %s", outfile, format, layout, scale, more)

	return system(cmd, dot)
}

func system(c string, dot string) error {
	cmd := exec.Command(`/bin/sh`, `-c`, c)
	cmd.Stdin = strings.NewReader(dot)
	return cmd.Run()

}
