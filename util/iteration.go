package util

func ForEach(in chan interface{}, f func()) {
	for range in {
		f()
	}
}
