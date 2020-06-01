package dataBase

type Status int

const (
	Pending Status = iota
	Rejected
	NotSelected
	Selected
	Started
	Waiting
	Accepted
	NotAccepted
	Expired
)

func (s Status) String() string {
	return [...]string{"Pending", "Rejected", "Not selected", "Selected", "Started",
		"Waiting", "Accepted", "Not accepted", "Expired"}[s]
}

func (s Status) IsValid() bool {
	switch s {
	case Pending, NotSelected, Selected, Started, Waiting, NotAccepted, Accepted, Expired, Rejected:
		return true
	}
	return false
}
