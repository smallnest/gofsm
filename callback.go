package fsm

// EventProcessor defines OnExit, Action and OnEnter actions.
type EventProcessor interface {
	// OnExit Action handles exiting a state
	OnExit(fromState string, args []interface{})
	// Action is used to handle transitions
	Action(action string, fromState string, toState string, args []interface{}) error
	// OnActionFailure failed to execute the Action
	OnActionFailure(action string, fromState string, toState string, args []interface{}, err error)
	// OnExit Action handles entering a state
	OnEnter(toState string, args []interface{})
}

// DefaultDelegate is a default delegate.
// it splits processing of actions into three actions: OnExit, Action and OnEnter.
type DefaultDelegate struct {
	P EventProcessor
}

// HandleEvent implements Delegate interface and split HandleEvent into three actions.
func (dd *DefaultDelegate) HandleEvent(action string, fromState string, toState string, args []interface{}) error {
	if fromState != toState {
		dd.P.OnExit(fromState, args)
	}

	err := dd.P.Action(action, fromState, toState, args)
	if err != nil {
		dd.P.OnActionFailure(action, fromState, toState, args, err)
		return err
	}

	if fromState != toState {
		dd.P.OnEnter(toState, args)
	}

	return nil
}
