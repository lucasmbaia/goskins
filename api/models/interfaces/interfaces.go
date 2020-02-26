package interfaces

type Models interface {
	Get(interface{}) (interface{}, error)
	Post(interface{}) (bool, error)
}
