package object

type Object interface {
	node()
}

type Int struct {
	Value int32
}

func (o Int) node() {}

type Float struct {
	Value float64
}

func (o Float) node() {}

type String struct {
	Value string
}

func (o String) node() {}

type Bool struct {
	Value bool
}

func (o Bool) node() {}
