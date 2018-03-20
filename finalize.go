package tripanel

var finalizations []func()

func OnQuit(f func()) {
	finalizations = append(finalizations, f)
}

func doquits() {
	for _, f := range finalizations {
		f()
	}
}
