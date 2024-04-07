package a

func call() int { return 0 }

var i = call() // want "call must not"

type X struct{}

func (X) call() int { return 0 }

var x = X{}.call() // want "X.call must not"

var z, y = call(), X{}.call() // want "call must not" "X.call must not"

var (
	a = call()     // want "call must not"
	b = X{}.call() // want "X.call must not"
)

func (*X) Call() int { return 0 }

var p = (&X{}).Call() // want `\*X.Call must not`
