package process

type StringMap map[string]string

type Event struct {
	Origin string
	Data   StringMap
}
