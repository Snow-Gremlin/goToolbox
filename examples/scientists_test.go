package examples

import (
	"cmp"
	_ "embed"
	"fmt"
	"strconv"

	"goToolbox/collections"
	"goToolbox/collections/enumerator"
	"goToolbox/collections/set"
	"goToolbox/collections/sortedDictionary"
	"goToolbox/collections/tuple4"
)

//go:embed scientists_data.md
var mdFile string

type entry collections.Tuple4[string, string, int, collections.Set[string]]

func newEntry(parts []string) entry {
	firstName := parts[0]
	lastName := parts[1]
	birthYear, err := strconv.Atoi(parts[2])
	if err != nil {
		panic(err)
	}
	focus := set.From(enumerator.Split(parts[3], `,`).Trim())
	return tuple4.New(firstName, lastName, birthYear, focus)
}

func getEntryEnumerator() collections.Enumerator[entry] {
	lines := enumerator.Lines(mdFile). // Breaks the text into lines
						Skip(4).  // Skips over the first 4 lines, the header, of the file.
						NotZero() // Skip blank lines

	allParts := enumerator.Select(lines, func(line string) []string {
		return enumerator.Split(line, `|`). // Split line into parts
							Trim().   // Trim space off of each part
							ToSlice() // Put the parts into a slice
	})

	return enumerator.Select(allParts, newEntry)
}

func sortByYear(a, b entry) int {
	c := cmp.Compare(a.Value3(), b.Value3()) // by year
	if c == 0 {
		c = cmp.Compare(a.Value2(), b.Value2()) // then by last name
		if c == 0 {
			c = cmp.Compare(a.Value1(), b.Value1()) // then by first name
		}
	}
	return c
}

func sortByLastName(a, b entry) int {
	c := cmp.Compare(a.Value2(), b.Value2()) // by last name
	if c == 0 {
		c = cmp.Compare(a.Value1(), b.Value1()) // then by first name
		if c == 0 {
			c = cmp.Compare(a.Value3(), b.Value3()) // then by year
		}
	}
	return c
}

func isPhysicists(a entry) bool {
	return a.Value4().Contains(`Physics`)
}

func inMedicine(a entry) bool {
	return a.Value4().Contains(`Medicine`)
}

func printEntries(entries collections.Enumerator[entry]) {
	enumerator.Indexed(entries).Foreach(func(t collections.Tuple2[int, entry]) {
		index := t.Value1()
		entry := t.Value2()
		fmt.Printf("\t%d. %s, %s (%d)\n", index+1, entry.Value2(), entry.Value1(), entry.Value3())
	})
}

func Example_isPhysicists_sortByYear() {
	entries := getEntryEnumerator().
		Where(isPhysicists).
		Sort(sortByYear)

	fmt.Println(`Physicists:`)
	printEntries(entries)

	// Output: Physicists:
	//	1. du Châtelet, Émilie (1706)
	//	2. Bassi, Laura (1711)
	//	3. Germain, Sophie (1776)
	//	4. Meitner, Lise (1878)
	//	5. Wu, Chien-Shiung (1912)
	//	6. Ride, Sally (1951)
	//	7. Profet, Margaret (1958)
}

func Example_inMedicine_sortByLastName() {
	entries := getEntryEnumerator().
		Where(inMedicine).
		Sort(sortByLastName)

	fmt.Println(`Physicians:`)
	printEntries(entries)

	// Output: Physicians:
	//	1. Anderson, Elizabeth (1836)
	//	2. Apgar, Virginia (1909)
	//	3. Barre-Sinoussi, Francoise (1947)
	//	4. Barton, Clara (1821)
	//	5. Bath, Patricia (1942)
	//	6. Blackwell, Elizabeth (1821)
	//	7. Cori, Gerty (1896)
	//	8. Elion, Gertrude (1918)
	//	9. Hamilton, Alice (1869)
	//	10. McClintock, Barbara (1902)
	//	11. Moser, May-Britt (1963)
	//	12. Nightingale, Florence (1820)
	//	13. Novello, Antonia (1944)
	//	14. Sabin, Florence (1871)
	//	15. Sanger, Margaret (1879)
	//	16. Stevenson, Sarah (1841)
	//	17. Taussig, Helen (1898)
	//	18. Yalow, Rosalyn (1921)
}

func getFocusCounts() collections.Dictionary[string, int] {
	focuses := sortedDictionary.New[string, int]()
	enumerator.Expand(getEntryEnumerator(), func(e entry) collections.Iterable[string] {
		return e.Value4().Enumerate().Iterate
	}).Foreach(func(focus string) {
		count, _ := focuses.TryGet(focus)
		focuses.Add(focus, count+1)
	})
	return focuses
}

func Example_getFocusCounts() {
	focuses := getFocusCounts()
	focuses.RemoveIf(func(focus string) bool {
		return focuses.Get(focus) <= 1
	})

	fmt.Println(`Focuses:`)
	fmt.Println(focuses.String())

	// Output: Focuses:
	// Agriculture:       2
	// Anthropology:      4
	// Astronomy:         5
	// Bacteriology:      3
	// Biology:           11
	// Calculus:          2
	// Chemistry:         4
	// Computer Science:  4
	// Education:         18
	// Entomology:        2
	// Environmentalist:  3
	// Genetics:          6
	// Geology:           3
	// Mathematics:       15
	// Medicine:          18
	// Microbiology:      4
	// Molecular Biology: 2
	// Neurology:         2
	// Nuclear Physics:   5
	// Nursing:           3
	// Paleontology:      3
	// Philosophy:        2
	// Physics:           7
	// Physiology:        8
	// Primatology:       3
	// Virology:          2
}
