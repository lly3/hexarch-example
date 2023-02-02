package domain

import (
	"errors"
	"fmt"
)

type State string

const (
	StateWaitingForPayment      State = "waiting-for-payment"
	StateWaitingForConfirmation State = "waiting-for-confirmation"
	StateWaitingForCooking      State = "waiting-for-cooking"
	StateDelivering             State = "delivering"
	StateComplete               State = "completed"
	StateCancel                 State = "cancelled"
)

type Order struct {
	Id           string
	ClientId     string
	RestaurantId string
	State        State
}

func (o *Order) UpdateOrderStatus(state State) error {
	if o.State == StateComplete {
		return errors.New("This order is already completed")
	}
	if o.State == StateCancel {
		return errors.New("This order is already cancelled")
	}

	if o.State == StateWaitingForPayment {
		if state == StateWaitingForConfirmation {
			o.State = state
			return nil
		}
		if state == StateCancel {
			o.State = state
			return nil
		}
	}
	if o.State == StateWaitingForConfirmation {
		if state == StateWaitingForCooking {
			o.State = state
			return nil
		}
		if state == StateCancel {
			o.State = state
			return nil
		}
	}
	if o.State == StateWaitingForCooking {
		if state == StateDelivering {
			o.State = state
			return nil
		}
		if state == StateCancel {
			o.State = state
			return nil
		}
	}
	if o.State == StateDelivering {
		if state == StateComplete {
			o.State = state
			return nil
		}
		if state == StateCancel {
			o.State = state
			return nil
		}
	}

	return fmt.Errorf("Invalid state from: %s to: %s", o.State, state)
}
