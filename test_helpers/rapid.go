package th

import "pgregory.net/rapid"

func NonEmptyString() *rapid.Generator[string] {
	return rapid.StringOfN(
		rapid.RuneFrom([]rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")), // not too crazy string for now
		1, -1, -1,
	)
}
