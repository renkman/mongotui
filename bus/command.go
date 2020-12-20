package bus

type Command struct {
	name       string
	getMessage func() string
}

type Listener interface {
}

type Foo struct {
}
