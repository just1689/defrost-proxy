package util

func IntToInt64P(i int) *int32 {
	x := int32(i)
	return &x
}
