package dataBase

import "errors"

type Status int

const (
	Pending Status = iota
	NotSelected
	Selected
	Started
	Waiting
	NotAccepted
	Accepted
	Expired
)

func (s Status) String() string {
	return [...]string{"Pending", "Not selected", "Selected", "Started",
		"Waiting", "Not accepted", "Accepted", "Expired"}[s]
}

func (s Status) IsValid() error {
	switch s {
	case Pending, NotSelected, Selected, Started, Waiting, NotAccepted, Accepted, Expired:
		return nil
	}
	return errors.New(WrongValueError.Error() + `status`)
}
