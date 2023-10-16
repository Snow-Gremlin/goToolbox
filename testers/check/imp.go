package check

import (
	"goToolbox/collections"
	"goToolbox/collections/readonlyVariantList"
	"goToolbox/testers"
	"goToolbox/utils"
)

type trialHandle[T any] func(b *testee, actual T)

type checkImp[T any] struct {
	b     *testee
	trial trialHandle[T]
}

func newCheck[T any](t testers.Tester, trial trialHandle[T]) *checkImp[T] {
	return &checkImp[T]{
		b:     newTestee(t),
		trial: trial,
	}
}

func newPred[T any](t testers.Tester, p collections.Predicate[T], action string) *checkImp[T] {
	return newCheck(t, func(b *testee, actual T) {
		if !p(actual) {
			b.Should(action)
		}
	})
}

func newLen[T any](t testers.Tester, p collections.Predicate[int], action string) *checkImp[T] {
	return newCheck(t, func(b *testee, actual T) {
		if length, ok := utils.Length(actual); !ok {
			b.Should(`be a type that has length`)
		} else if !p(length) {
			b.With(`Actual Length`, length).
				Should(action)
		}
	})
}

func (c *checkImp[T]) copyAndAdd(handle func(c2 *checkImp[T])) *checkImp[T] {
	if c == nil {
		return nil
	}

	c2 := &checkImp[T]{
		b:     c.b.Copy(),
		trial: c.trial,
	}
	handle(c2)
	return c2
}

func (c *checkImp[T]) With(key string, args ...any) testers.Check[T] {
	return c.copyAndAdd(func(c2 *checkImp[T]) {
		c2.b.With(key, args...)
	})
}

func (c *checkImp[T]) Withf(key, format string, args ...any) testers.Check[T] {
	return c.copyAndAdd(func(c2 *checkImp[T]) {
		c2.b.Withf(key, format, args...)
	})
}

func (c *checkImp[T]) WithType(key string, valueForType any) testers.Check[T] {
	return c.copyAndAdd(func(c2 *checkImp[T]) {
		c2.b.WithType(key, valueForType)
	})
}

func (c *checkImp[T]) WithValue(key string, value any) testers.Check[T] {
	return c.copyAndAdd(func(c2 *checkImp[T]) {
		c2.b.WithValue(key, value)
	})
}

func (c *checkImp[T]) Name(name string) testers.Check[T] {
	return c.With(`Name`, name)
}

func (c *checkImp[T]) AsText() testers.Check[T] {
	return c.copyAndAdd(func(c2 *checkImp[T]) {
		c2.b.textHint = true
	})
}

func (c *checkImp[T]) AsHex() testers.Check[T] {
	return c.copyAndAdd(func(c2 *checkImp[T]) {
		c2.b.intHint = 16
	})
}

func (c *checkImp[T]) AsOct() testers.Check[T] {
	return c.copyAndAdd(func(c2 *checkImp[T]) {
		c2.b.intHint = 8
	})
}

func (c *checkImp[T]) AsBin() testers.Check[T] {
	return c.copyAndAdd(func(c2 *checkImp[T]) {
		c2.b.intHint = 2
	})
}

func (c *checkImp[T]) TimeAs(timeFormat string) testers.Check[T] {
	return c.copyAndAdd(func(c2 *checkImp[T]) {
		c2.b.timeHint = timeFormat
	})
}

func (c *checkImp[T]) Required() testers.Check[T] {
	return c.copyAndAdd(func(c2 *checkImp[T]) {
		c2.b.Required()
	})
}

func (c *checkImp[T]) Require(actual T) testers.Check[T] {
	return c.Required().Assert(actual)
}

func (c *checkImp[T]) Assert(actual T) (pc testers.Check[T]) {
	if c == nil {
		return c
	}

	defer handlePanic(c.b.t, &pc)
	getHelper(c.b.t)()

	b2 := c.b.Copy().
		WithValue(`Actual Value`, actual).
		WithType(`Actual Type`, actual)
	c.trial(b2, actual)
	b2.Finish()
	return c
}

func (c *checkImp[T]) AssertAll(actual any) (pc testers.Check[T]) {
	if c == nil {
		return c
	}

	defer handlePanic(c.b.t, &pc)
	getHelper(c.b.t)()

	actV := readonlyVariantList.Cast[T](readonlyVariantList.Wrap(actual))
	it := actV.Enumerate().Iterate()
	index := 0
	for it.Next() {
		value := it.Current()
		b2 := c.b.Copy().
			WithValue(`Actual Collection`, actual).
			WithValue(`Actual Value`, value).
			WithType(`Actual Type`, value).
			WithValue(`Value Index`, index)
		c.trial(b2, value)
		b2.Finish()
		index++
	}
	return c
}

func (c *checkImp[T]) RequireAll(actual any) testers.Check[T] {
	return c.Required().AssertAll(actual)
}

func (c *checkImp[T]) Panic(handle func()) (pc testers.Check[T]) {
	if c == nil {
		return nil
	}

	defer handlePanic(c.b.t, &pc)
	getHelper(c.b.t)()

	recovered := func() (r any) {
		defer func() { r = recover() }()
		handle()
		return nil
	}()

	b2 := c.b.Copy()
	if recovered == nil {
		b2.Should(`panic from given function`).
			Finish()
		return c
	}

	actual, ok := recovered.(T)
	if !ok {
		b2.WithValue(`Panicked Value`, recovered).
			WithType(`Panicked Type`, recovered).
			Should(`be a panic of the expected type`).
			Finish()
		return c
	}

	b2 = c.b.Copy().
		WithValue(`Panicked Value`, actual).
		WithType(`Panicked Type`, actual)
	c.trial(b2, actual)
	b2.Finish()
	return c
}

func (c *checkImp[T]) setTextHint(value any) *checkImp[T] {
	if _, ok := value.(string); ok {
		c.b.textHint = true
	}
	return c
}
