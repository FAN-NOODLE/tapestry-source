package modes

type Mode struct {
	Params interface{}
}

func (m *Mode) SetParams (params interface{})  {
	m.Params = params
}