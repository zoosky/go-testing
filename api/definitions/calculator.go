package definitions

// CalculatorRequest represents a generic calculator operation request
type CalculatorRequest struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
}

// CalculatorResponse represents a generic calculator operation response
type CalculatorResponse struct {
	Result float64 `json:"result"`
}