//参考了 https://github.com/elimisteve/fsm 的实现

package fsm

import (
	"fmt"
	"os/exec"
	"strings"
)

// Transition 代表一个转换。所有的数据都是以字符串字面量表示的，这样可以简化状态机的实现。
type Transition struct {
	From   string
	Event  string
	To     string
	Action string
}

// Delegate 用来执行特定的action. 它可以根据action的字符串字面量以及补充参数执行特定的动作。
// Action的执行可以是同步的，也可以是异步的，在同步的情况下会阻塞对事件的处理。
// 在处理的时候处理函数可以检查对象的状态是否和fromState一致，如果不一致需要根据业务自行处理或者报错。
// fromState 和 toState 可以相同，这种情况下对象的状态不发生改变，只是需要处理事件而已。
type Delegate interface {
	// HandleEvent 处理事件和转换
	HandleEvent(action string, fromState string, toState string, args []interface{})
}

// StateMachine 用来代表状态机对象。当我们说状态的时候，肯定指的是某个对象的状态，我们要处理的就是触发这个对象的事件和这个对象的状态的改变。
// StateMachine 本身不保存对象的当前状态,所以它可以处理N个对象的状态转换，这也要求对象必须传入自己当前的状态。
// 因此这个StateMachine对象也可以称之为 “对象无关" 的状态机。
// 转换时执行的特定动作通过 delegate 实现。
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

// NewStateMachine 创建一个新的状态机
func NewStateMachine(delegate Delegate, transitions ...Transition) *StateMachine {
	return &StateMachine{delegate: delegate, transitions: transitions}
}

// Trigger 处理一个事件
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

// findTransMatching 根据输入(event)和当前状态找到转换对象。
func (m *StateMachine) findTransMatching(fromState string, event string) *Transition {
	for _, v := range m.transitions {
		if v.From == fromState && v.Event == event {
			return &v
		}
	}
	return nil
}

// Export 输出状态图到指定的文件
func (m *StateMachine) Export(outfile string) error {
	return m.ExportWithDetails(outfile, "png", "dot", "72", "-Gsize=10,5 -Gdpi=200")
}

// ExportWithDetails 输出状态图，提供更多的graphviz参数
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
