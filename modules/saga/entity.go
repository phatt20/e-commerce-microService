package saga

type SagaState string

const (
	SagaStarted   SagaState = "STARTED"
	SagaCompleted SagaState = "COMPLETED"
	SagaFailed    SagaState = "FAILED"
)

type SagaStep struct {
	Name      string
	Completed bool
	Data      []byte
}

type Saga struct {
	SagaID  string
	OrderID string
	State   SagaState
	Steps   []SagaStep
}
