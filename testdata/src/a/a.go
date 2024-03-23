package a

func call() int { return 0 }

var i = call() // want "x"

type X struct{}

func (X) call() int { return 0 }

var x = X{}.call() // want "x"

var z, y = call(), X{}.call() // want "x" "x"

var (
	a = call()     // want "x"
	b = X{}.call() // want "x"
)
