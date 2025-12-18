package jobqueue

type Status uint8

const (
	Pending Status = iota
	Processing
	Completed
	Failed
)

func (status Status) String() string {
	switch status {
	case Pending:
		return "Pending"
	case Processing:
		return "Processing"
	case Completed:
		return "Completed"
	case Failed:
		return "Failed"
	default:
		return ""
	}
}
