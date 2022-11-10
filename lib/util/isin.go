package util

func IsIn[E comparable](e E, in ...E) bool {
	for _, f := range in {
		if f == e {
			return true
		}
	}
	return false
}
