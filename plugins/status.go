package plugins

type ExecutionStatusName string

const (
	ExecutionStatusUnprocessed ExecutionStatusName = "unprocessed"
	ExecutionStatusCompleted   ExecutionStatusName = "completed"
	ExecutionStatusFailed      ExecutionStatusName = "failed"
)

type ExecutionStatus struct {
	Status  ExecutionStatusName `json:"status" yaml:"status"`
	Message string              `json:"message" yaml:"message"`
	// This is used for workflow execution to store the last completed step index
	LastCompletedStepIndex int `json:"last_completed_step_index" yaml:"last_completed_step_index"`
}

func (s *ExecutionStatus) GetStatus() ExecutionStatusName {
	if s == nil || s.Status == "" {
		return ExecutionStatusUnprocessed
	}
	return s.Status
}

func (s *ExecutionStatus) IsCompleted() bool {
	return s.Status == ExecutionStatusCompleted
}

func (s *ExecutionStatus) IsFailed() bool {
	return s.Status == ExecutionStatusFailed
}

func (s *ExecutionStatus) IsUnprocessed() bool {
	return s.Status == ExecutionStatusUnprocessed || s.Status == ""
}

func (s *ExecutionStatus) SetError(err error) {
	s.Status = ExecutionStatusFailed
	s.Message = err.Error()
}
