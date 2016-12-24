package fsm

// EventProcessor 接口定义了OnExit和OnEnter的hook.
type EventProcessor interface {
	OnExit(fromState string, args []interface{})
	Action(action string, fromState string, toState string, args []interface{})
	OnEnter(toState string, args []interface{})
}

// DefaultDelegate 在处理Event的时候，将事件的处理分成三步，提供了调用OnExit和OnEnter的hook的功能.
type DefaultDelegate struct {
	p EventProcessor
}

// HandleEvent 将处理事件分成三部步，退出前一个状态OnExit， 执行Action 和 进入下一个状态OnEnter.
func (dd *DefaultDelegate) HandleEvent(action string, fromState string, toState string, args []interface{}) {
	if fromState != toState {
		dd.p.OnExit(fromState, args)
	}

	dd.p.Action(action, fromState, toState, args)

	if fromState != toState {
		dd.p.OnEnter(toState, args)
	}
}
