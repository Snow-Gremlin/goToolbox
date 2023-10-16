package check

// TODO: Remove

/*
// hasItems checks a data structure for a given list of expected values.
// This will post errors if there are missing expected values.
func hasItemsWithPost[S any, T comparable](t testers.Tester, tag string, expected []T, handle func(S, T) bool) testers.Check[S] {
	count := len(expected)
	if count <= 0 {
		newPoster(t).Required().
			Withf(`Type`, `%T`, expected).
			PostAssert(`provide at least one expected ` + tag)
		return (*checkImp[S])(nil)
	}

	if count == 1 {
		return newCheck(t, func(p *poster, actual S) {
			if !handle(actual, expected[0]) {
				p.Withf(`Actual`, `%v`, actual).
					Withf(`Type`, `%T`, actual).
					Withf(`Expected`, `%v`, expected[0]).
					PostAssert(`have the expected ` + tag)
			}
		})
	}

	return newCheck(t, func(p *poster, actual S) {
		missing := []T{}
		for _, elem := range expected {
			if !handle(actual, elem) {
				missing = append(missing, elem)
			}
		}

		if len(missing) > 0 {
			p.Withf(`Actual`, `%v`, actual).
				Withf(`Type`, `%T`, actual).
				Withf(`Expected`, `%v`, expected).
				Withf(`Missing`, `%v`, missing).
				PostAssert(`have all the expected ` + tag + `s`)
		}
	})
}
*/
