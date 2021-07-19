package thread

func GoSafe(fn func()) {
	go RunSafe(fn)
}

func RunSafe(fn func()) {
	defer Recover()

	fn()
}
