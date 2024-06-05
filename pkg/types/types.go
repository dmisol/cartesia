package types

type Voice [192]float64
type Model string

type Request struct {
	ContextId string      `json:"context_id"`
	Data      RequestData `json:"data"`
}

type RequestData struct {
	Text  string `json:"transcript"`
	Model string `json:"model_id"`
	Voice Voice  `json:"voice"`
}
type Response struct {
	ContextId string `json:"context_id"`
	Done      bool   `json:"done"`

	Data         string  `json:"data,omitempty"`
	Length       int     `json:"length,omitempty"`
	SamplingRate int     `json:"sampling_rate,omitempty"`
	StepTime     float64 `json:"step_time,omitempty"`
}
