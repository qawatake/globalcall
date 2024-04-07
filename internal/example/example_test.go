package example_test

func Call() int { return 0 }

var i = Call() // <- ng

type X struct{}

func (X) Call() int { return 0 }

var j = X{}.Call() // <- ng

type Y struct{}

func (*Y) Call() int { return 0 }

var k = (&Y{}).Call() // <- ng
