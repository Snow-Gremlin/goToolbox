package diff

import (
	"strings"
	"testing"

	"goToolbox/collections/stack"
	"goToolbox/collections/tuple2"
	"goToolbox/differs"
	"goToolbox/differs/data"
	"goToolbox/differs/diff/internal/result"
	"goToolbox/differs/step"
	"goToolbox/utils"
)

var (
	exampleA = lines(
		`This part of the`,
		`document has stayed the`,
		`same from version to`,
		`version.  It shouldn't`,
		`be shown if it doesn't`,
		`change.  Otherwise, that`,
		`would not be helping to`,
		`compress the size of the`,
		`changes.`,
		``,
		`This paragraph contains`,
		`text that is outdated.`,
		`It will be deleted in the`,
		`near future.`,
		``,
		`It is important to spell`,
		`check this dokument. On`,
		`the other hand, a`,
		`misspelled word isn't`,
		`the end of the world.`,
		`Nothing in the rest of`,
		`this paragraph needs to`,
		`be changed. Things can`,
		`be added after it.`)

	exampleB = lines(
		`This is an important`,
		`notice! It should`,
		`therefore be located at`,
		`the beginning of this`,
		`document!`,
		``,
		`This part of the`,
		`document has stayed the`,
		`same from version to`,
		`version.  It shouldn't`,
		`be shown if it doesn't`,
		`change.  Otherwise, that`,
		`would not be helping to`,
		`compress anything.`,
		``,
		`It is important to spell`,
		`check this document. On`,
		`the other hand, a`,
		`misspelled word isn't`,
		`the end of the world.`,
		`Nothing in the rest of`,
		`this paragraph needs to`,
		`be changed. Things can`,
		`be added after it.`,
		``,
		`This paragraph contains`,
		`important new additions`,
		`to this document.`)

	billNyeA = `The most serious problem facing humankind is climate change. ` +
		`All of these people breathing and burning our atmosphere has led to an ` +
		`extraordinarily dangerous situation. I hope next generation will emerge ` +
		`and produce technology, regulations, and a worldview that enable as many ` +
		`of us as possible to live happy healthy lives.`

	billNyeB = `The meaning of life is pretty clear: Living things strive to ` +
		`pass their genes into the future. The claim that we would not have ` +
		`morals or ethics without religion is extraordinary. Animals in nature ` +
		`seem to behave in moral ways without organized religion.`

	hirschbergPlusMinus = lines(
		`+This is an important`,
		`+notice! It should`,
		`+therefore be located at`,
		`+the beginning of this`,
		`+document!`,
		`+`,
		` This part of the`,
		` document has stayed the`,
		` same from version to`,
		` version.  It shouldn't`,
		` be shown if it doesn't`,
		` change.  Otherwise, that`,
		` would not be helping to`,
		`-compress the size of the`,
		`-changes.`,
		`-`,
		`-This paragraph contains`,
		`-text that is outdated.`,
		`-It will be deleted in the`,
		`-near future.`,
		`+compress anything.`,
		` `,
		` It is important to spell`,
		`-check this dokument. On`,
		`+check this document. On`,
		` the other hand, a`,
		` misspelled word isn't`,
		` the end of the world.`,
		` Nothing in the rest of`,
		` this paragraph needs to`,
		` be changed. Things can`,
		` be added after it.`,
		`+`,
		`+This paragraph contains`,
		`+important new additions`,
		`+to this document.`)

	// wagner is different because of differences in which
	// equal Levenstein distance paths are preferences.
	wagnerPlusMinus = lines(
		`+This is an important`,
		`+notice! It should`,
		`+therefore be located at`,
		`+the beginning of this`,
		`+document!`,
		`+`,
		` This part of the`,
		` document has stayed the`,
		` same from version to`,
		` version.  It shouldn't`,
		` be shown if it doesn't`,
		` change.  Otherwise, that`,
		` would not be helping to`,
		`-compress the size of the`,
		`-changes.`,
		`+compress anything.`,
		` `,
		`-This paragraph contains`,
		`-text that is outdated.`,
		`-It will be deleted in the`,
		`-near future.`,
		`-`,
		` It is important to spell`,
		`-check this dokument. On`,
		`+check this document. On`,
		` the other hand, a`,
		` misspelled word isn't`,
		` the end of the world.`,
		` Nothing in the rest of`,
		` this paragraph needs to`,
		` be changed. Things can`,
		` be added after it.`,
		`+`,
		`+This paragraph contains`,
		`+important new additions`,
		`+to this document.`)
)

func Test_Diff_BadCounts(t *testing.T) {
	s := stack.With(
		tuple2.New(step.Equal, 2),
		tuple2.New(step.Added, 2),
		tuple2.New(step.Removed, 1))
	r := result.New(s.Enumerate(), s.Count(), 4, 3, 5, 2, 1)

	checkPanic(t, `must have the same number of values as the result was created with `+
		`{incorrect A count: 3, required A count: 4, required B count: 3}`,
		func() { EnumerateValues(r, []int{1, 2, 3}, []int{1, 2, 5}) })

	checkPanic(t, `must have the same number of values as the result was created with `+
		`{incorrect B count: 4, required A count: 4, required B count: 3}`,
		func() { EnumerateValues(r, []int{1, 2, 3, 4}, []int{1, 2, 4, 5}) })

	checkPanic(t, `must have the same number of values as the result was created with `+
		`{incorrect A count: 3, incorrect B count: 4, required A count: 4, required B count: 3}`,
		func() { EnumerateValues(r, []int{1, 2, 3}, []int{1, 2, 4, 5}) })
}

func Test_Diff_Basics(t *testing.T) {
	checkLP(t, "A", "A", "=1")
	checkLP(t, "A", "B", "-1 +1")
	checkLP(t, "A", "AB", "=1 +1")
	checkLP(t, "A", "BA", "+1 =1")
	checkLP(t, "AB", "A", "=1 -1")
	checkLP(t, "BA", "A", "-1 =1")
	checkLP(t, "kitten", "sitting", "-1 +1 =3 -1 +1 =1 +1")
	checkLP(t, "saturday", "sunday", "=1 -2 =1 -1 +1 =3")
	checkLP(t, "satxrday", "sunday", "=1 -4 +2 =3")
	checkLP(t, "ABC", "ADB", "=1 +1 =1 -1")
}

func Test_Diff_Words(t *testing.T) {
	wordsA := strings.Split(billNyeA, ` `)
	wordsB := strings.Split(billNyeB, ` `)
	path := Default().Diff(data.Strings(wordsA, wordsB))
	exp := `=1 -9 +1 =1 -9 +7 =1 -17 +8 =1 -7 +15 =1 -4 +7`
	if result := utils.String(path); exp != result {
		t.Error("Diff returned unexpected result:",
			"\n   Input A:  ", wordsA,
			"\n   Input B:  ", wordsB,
			"\n   Expected: ", exp,
			"\n   Result:   ", result)
	}
}

func Test_Diff_PlusMinus_Parts(t *testing.T) {
	checkPlusMinus(t, ",",
		"cat,dog,pig",
		"cat,horse,dog",
		" cat,+horse, dog,-pig")
	checkPlusMinus(t, ",",
		"mike,ted,mark,jim",
		"ted,mark,bob,bill",
		"-mike, ted, mark,-jim,+bob,+bill")
	checkPlusMinus(t, ",",
		"k,i,t,t,e,n",
		"s,i,t,t,i,n,g",
		"-k,+s, i, t, t,-e,+i, n,+g")
	checkPlusMinus(t, ",",
		"s,a,t,u,r,d,a,y",
		"s,u,n,d,a,y",
		" s,-a,-t, u,-r,+n, d, a, y")
	checkPlusMinus(t, ",",
		"s,a,t,x,r,d,a,y",
		"s,u,n,d,a,y",
		" s,-a,-t,-x,-r,+u,+n, d, a, y")
	checkPlusMinus(t, ",",
		"func A() int,{,return 10,},,func C() int,{,return 12,}",
		"func A() int,{,return 10,},,func B() int,{,return 11,},,func C() int,{,return 12,}",
		" func A() int, {, return 10, }, ,+func B() int,+{,+return 11,+},+, func C() int, {, return 12, }")
}

func Test_Diff_Inline_Parts(t *testing.T) {
	checkInline(t, ",",
		"cat,dog,pig",
		"cat,horse,dog",
		" cat,+horse, dog,-pig")
	checkInline(t, ",",
		"mike,ted,mark,jim",
		"ted,mark,bob,bill",
		"-mike, ted_mark,-jim,+bob_bill")
	checkInline(t, ",",
		"k,i,t,t,e,n",
		"s,i,t,t,i,n,g",
		"-k,+s, i_t_t,-e,+i, n,+g")
	checkInline(t, ",",
		"s,a,t,u,r,d,a,y",
		"s,u,n,d,a,y",
		" s,-a_t, u,-r,+n, d_a_y")
	checkInline(t, ",",
		"s,a,t,x,r,d,a,y",
		"s,u,n,d,a,y",
		" s,-a_t_x_r,+u_n, d_a_y")
	checkInline(t, ",",
		"func A() int,{,return 10,},,func C() int,{,return 12,}",
		"func A() int,{,return 10,},,func B() int,{,return 11,},,func C() int,{,return 12,}",
		" func A() int_{_return 10_}_,+func B() int_{_return 11_}_, func C() int_{_return 12_}")
}

func Test_Diff_Merge_Lines(t *testing.T) {
	checkSlices(t, Default().Merge(exampleA, exampleB), lines(
		`<<<<<<<<`,
		`========`,
		`This is an important`,
		`notice! It should`,
		`therefore be located at`,
		`the beginning of this`,
		`document!`,
		``,
		`>>>>>>>>`,
		`This part of the`,
		`document has stayed the`,
		`same from version to`,
		`version.  It shouldn't`,
		`be shown if it doesn't`,
		`change.  Otherwise, that`,
		`would not be helping to`,
		`<<<<<<<<`,
		`compress the size of the`,
		`changes.`,
		``,
		`This paragraph contains`,
		`text that is outdated.`,
		`It will be deleted in the`,
		`near future.`,
		`========`,
		`compress anything.`,
		`>>>>>>>>`,
		``,
		`It is important to spell`,
		`<<<<<<<<`,
		`check this dokument. On`,
		`========`,
		`check this document. On`,
		`>>>>>>>>`,
		`the other hand, a`,
		`misspelled word isn't`,
		`the end of the world.`,
		`Nothing in the rest of`,
		`this paragraph needs to`,
		`be changed. Things can`,
		`be added after it.`,
		`<<<<<<<<`,
		`========`,
		``,
		`This paragraph contains`,
		`important new additions`,
		`to this document.`,
		`>>>>>>>>`))
}

func Test_Diff_Merge_MoreCases(t *testing.T) {
	checkSlices(t, Default().Merge(lines(
		`sameA`,
		`removedA`,
		`sameB`,
		`sameC`,
		`removedC`,
		`sameD`,
	), lines(
		`sameA`,
		`sameB`,
		`AddedB`,
		`sameC`,
		`AddedC`,
		`sameD`,
	)), lines(
		`sameA`,
		`<<<<<<<<`,
		`removedA`,
		`========`,
		`>>>>>>>>`,
		`sameB`,
		`<<<<<<<<`,
		`========`,
		`AddedB`,
		`>>>>>>>>`,
		`sameC`,
		`<<<<<<<<`,
		`removedC`,
		`========`,
		`AddedC`,
		`>>>>>>>>`,
		`sameD`,
	))

	checkSlices(t, Default().Merge(lines(
		`sameA`,
		`removedA`,
	), lines(
		`sameA`,
	)), lines(
		`sameA`,
		`<<<<<<<<`,
		`removedA`,
		`========`,
		`>>>>>>>>`,
	))

	checkSlices(t, Default().Merge(lines(
		`sameA`,
	), lines(
		`sameA`,
		`addedA`,
	)), lines(
		`sameA`,
		`<<<<<<<<`,
		`========`,
		`addedA`,
		`>>>>>>>>`,
	))
}

func Test_Diff_Merge_EdgeCases(t *testing.T) {
	checkSlices(t, Default().Merge(lines(
		`sameA`,
		`removedA`,
	), lines(
		`sameA`,
		`AddedA`,
	)), lines(
		`sameA`,
		`<<<<<<<<`,
		`removedA`,
		`========`,
		`AddedA`,
		`>>>>>>>>`,
	))

	// Normally remove is first but check if added is first.
	s := stack.With(
		tuple2.New(step.Equal, 1),
		tuple2.New(step.Added, 1),
		tuple2.New(step.Removed, 1))
	checkSlices(t, Merge(result.New(s.Enumerate(), s.Count(), 2, 2, 3, 1, 1), lines(
		`sameA`,
		`removedA`,
	), lines(
		`sameA`,
		`AddedA`,
	)), lines(
		`sameA`,
		`<<<<<<<<`,
		`========`,
		`AddedA`,
		`>>>>>>>>`,
		`<<<<<<<<`,
		`removedA`,
		`========`,
		`>>>>>>>>`,
	))
}

func Test_Diff_PlusMinus_Lines(t *testing.T) {
	checkSlices(t, Default().PlusMinus(exampleA, exampleB), hirschbergPlusMinus)

	checkSlices(t, Hirschberg(-1, false).PlusMinus(exampleA, exampleB), hirschbergPlusMinus)
	checkSlices(t, Hirschberg(-1, true).PlusMinus(exampleA, exampleB), hirschbergPlusMinus)

	checkSlices(t, Hybrid(-1, false, -1).PlusMinus(exampleA, exampleB), hirschbergPlusMinus)
	checkSlices(t, Hybrid(-1, true, -1).PlusMinus(exampleA, exampleB), hirschbergPlusMinus)

	checkSlices(t, Wagner(-1).PlusMinus(exampleA, exampleB), wagnerPlusMinus)
}

func lines(ln ...string) []string {
	return ln
}

func checkLP(t *testing.T, a, b, exp string) {
	checkLPDiff(t, Wagner(-1), `Wagner`, a, b, exp)
	checkLPDiff(t, Hirschberg(-1, true), `Default`, a, b, exp)
	checkLPDiff(t, Default(), `Default`, a, b, exp)
}

func checkLPDiff(t *testing.T, diff differs.Diff, alg, a, b, exp string) {
	path := diff.Diff(data.Chars(a, b))
	if result := utils.String(path); exp != result {
		t.Error("Diff returned unexpected result:",
			"\n   Algorithm: ", alg,
			"\n   Input A:   ", a,
			"\n   Input B:   ", b,
			"\n   Expected:  ", exp,
			"\n   Result:    ", result)
	}
}

func checkPlusMinus(t *testing.T, sep, a, b, exp string) {
	aParts := strings.Split(a, sep)
	bParts := strings.Split(b, sep)
	resultParts := Default().PlusMinus(aParts, bParts)
	result := strings.Join(resultParts, sep)
	if exp != result {
		t.Error("PlusMinus returned unexpected result:",
			"\n   Input A:  ", a,
			"\n   Input B:  ", b,
			"\n   Expected: ", exp,
			"\n   Result:   ", result)
	}
}

func checkInline(t *testing.T, sep, a, b, exp string) {
	aParts := strings.Split(a, sep)
	bParts := strings.Split(b, sep)
	r := Default().Diff(data.Strings(aParts, bParts))
	resultParts := Inline(r, aParts, bParts, `_`)
	result := strings.Join(resultParts, sep)
	if exp != result {
		t.Error("Inline returned unexpected result:",
			"\n   Input A:  ", a,
			"\n   Input B:  ", b,
			"\n   Expected: ", exp,
			"\n   Result:   ", result)
	}
}

func checkSlices(t *testing.T, result, exp []string) {
	resultStr := strings.Join(result, "\n")
	expStr := strings.Join(exp, "\n")
	if expStr != resultStr {
		t.Error("Unexpected result:",
			"\n   Expected: ", expStr,
			"\n   Result:   ", resultStr)
	}
}

func checkPanic(t *testing.T, exp string, handle func()) {
	result := func() (p string) {
		defer func() {
			p = utils.String(recover())
		}()
		handle()
		return ``
	}()
	if exp != result {
		t.Error("unexpected panic was gotten:",
			"\n   Expected: ", exp,
			"\n   Result:   ", result)
	}
}
