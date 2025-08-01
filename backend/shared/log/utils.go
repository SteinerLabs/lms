package log

// shallowCopy performs a shallow copy of the given map[string]any and
// returns a new map with the copied key-value pairs. If the given map is empty,
// it returns an empty map.
//
// Example:
//
//	m := map[string]any{"foo": 1, "bar": 2}
//	cp := shallowCopy(m) // cp is a new map containing the key-value pairs of m
//
//	empty := shallowCopy(map[string]any{}) // empty is an empty map
//
// Note: The function does not perform a deep copy, so any nested maps or
// reference types will be shared between the original and copied map.
func shallowCopy(m map[string]any) map[string]any {
	cp := make(map[string]any, len(m))

	if len(m) > 0 {
		for k, v := range m {
			cp[k] = v
		}
	}

	return cp
}

// merge performs a merge operation on two maps and returns a new map with the merged key-value pairs.
// It takes two parameters, m0 and m1, which are maps with string keys and values of any type.
// The function starts by making a shallow copy of the m0 map using the shallowCopy function.
// If m1 is not empty, it iterates over each key-value pair in m1 and adds them to the copied map.
// Finally, the function returns the copied map.
//
// Example:
//
//	m0 := map[string]any{"foo": 1, "bar": 2}
//	m1 := map[string]any{"baz": 3}
//	result := merge(m0, m1) // result is a new map containing {"foo": 1, "bar": 2, "baz": 3}
//
// Note: The function does not modify the original maps and returns a new map with the merged key-value pairs.
func merge(m0 map[string]any, m1 map[string]any) map[string]any {
	cp := shallowCopy(m0)

	if len(m1) > 0 {
		for k, v := range m1 {
			cp[k] = v
		}
	}

	return cp
}
