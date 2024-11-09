package enums

type TaskStatus uint8

const (
	ToDo TaskStatus = iota
	InProgress
	Done
)

func TaskStatusFromString(taskString string) TaskStatus {
	switch taskString {
	case ToDo.String():
		return ToDo
	case InProgress.String():
		return InProgress
	case Done.String():
		return Done
	default:
		return ToDo
	}
}
