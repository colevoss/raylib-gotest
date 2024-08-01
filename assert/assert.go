package assert

func Assert(assertion bool, msg string) {
	if !assertion {
		panic(msg)
	}
}
