package a

func call() int { return 0 }

var i = call() // want "y"

type X struct{}

func (X) call() int { return 0 }

var x = X{}.call() // want "x"

var z, y = call(), X{}.call() // want "y" "x"

var (
	a = call()     // want "y"
	b = X{}.call() // want "x"
)

func (*X) Call() int { return 0 }

var p = (&X{}).Call() // want "x"
