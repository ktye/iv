package primitives

/* TODO remove

// Uptype tries to apply the given function to the value R.
// If it returns false, R is uptyped to the next higher numeric type.
func uptype(a *apl.Apl, R apl.Value, name string, apply func(n apl.Value) (apl.Value, bool)) (apl.Value, error) {
	// Try directly. This works even for types which are no numbers.
	if v, ok := apply(R); ok {
		return v, nil
	}

	n, ok := R.(apl.Number)
	if ok == false {
		return nil, fmt.Errorf("%s: expected a number: %T", name, R)
	}
	num, ok := a.Tower.Numbers[reflect.TypeOf(n)]
	if ok == false {
		return nil, fmt.Errorf("%s: unknown numeric type: %T", name, R)
	}
	for i := num.Class; i < a.Tower.Numbers; i++ {
		if res, ok := apply(n); ok {
			return res, nil
		}
		n, ok = num.Uptype(n)
		if ok == false {
			break
		}
		num = a.Tower.Numbers[reflect.TypeOf(n)]
	}
	return nil, fmt.Errorf("%s: not supported for %T", name, R)
}
*/
