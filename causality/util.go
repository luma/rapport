package causality

func dotsEql(d1 Dots, d2 Dots) bool {
	lenD1 := len(d1)
	lenD2 := len(d2)

	if lenD1 != lenD2 {
		return false
	}

	if lenD1 < lenD2 {
		for k, v := range d1 {
			if d2[k] != v {
				return false
			}
		}
	}

	for k, v := range d2 {
		if d1[k] != v {
			return false
		}
	}

	return true
}
