# TODO

These are ideas for new things that could be added to the toolbox
and things which need to be updated or improved.
This also contains ideas which need to be investigated.

## New Features

- [ ] `RingQueue`: Implements `collections.Queue` interface.
      Using two indices on a slice. Slice grows if out of room.

- [ ] `SliceStack`: Implements `collections.Stack` interface.
      Using an index on a slice. Slice grows if out of room.

- [ ] `SortedList`: Implements a new `collections.SortedList` interface that
      extends `collections.ReadonlyList`. A sorted list won't have the same
      insert at index methods as an unsorted list.

- [ ] `SortedSet`: Implements `collections.Set` interface, but uses a comparator
      function to sort the value. The `SortedSet` may leverage the
      `SortedDictionary` as storage.

- [ ] `FixedList`: Implements a new `collection.FixedList` interface that
      extends `collections.ReadonlyList`. Obviously the list must be defined
      with a size and may not modify that length after creation.

- [ ] Investigate some kind of B-tree for a different implementation of a list,
      sorted set, and other collections. The list would work like a variable
      list/tree of pointers to arrays such that each array is fixed length
      thus giving capacity to the list but the list is variable length by adding
      more and removing unused fixed arrays. This should minimize the amount of
      values that have to be copied anytime the list is changed at the expense
      of using extra memory for the extra capacity in each node.

- [ ] `Interpolator`: This takes a `collections.Getter[int]` with
      `collections.Countable`. It will provide access the data using a floating
      point index. The values in the list should contain floating points or
      integers, otherwise reading it would needs some kind of weighted value
      resolver (e.g. What is half way between `"cat"` and `"dog"`?).
      This implements linear interpolation, splines, and other methods of
      interpolating the data using one or more values from the underlying list.
      Maybe even implement ease-in, overshoot, ease-out, and other smoothers.

- [ ] `PriorityQueue`: Implement a Fibonacci priority queue.

- [ ] Add predicate and check for "multiple of" which uses a modulus to
      determine if an integer or floating point value is a multiple of
      some value.

- [ ] Add format hints, like `Hex()`, `Bin()`, etc to `Check`.
      These give hints to how to format the expected and actual values.
      Also think about adding key groups where keys contain a `;`,
      like `key (';' <key> )*`, group such that a tree of keys is built
      where each leaf has a value on it.

## New Enumerator Functions

- [ ] `GroupBy(e Enumerator[T], Selector[T, TKey]) Enumerator[Tuple[TKey, []T]`:
      groups the values together based on the keys then enumerates the groups.

- [ ] `TakeLast(count int) Enumerator[T]`: takes the last count number of values
      and returns them.

## New List Functions

- [ ] `SubRange(index, count int) List[T]`: gets a copy of a subRange of a list.
      Methods like this could also use a new list type that is a windowed
      version of the full list, then as the window is modified the parent list
      is modified and the window updates to changes to the parent list.

- [ ] `TakeRange(index, count int) List[T]`: gets and removes a subRange of a list.

- [ ] `Front(count int) List[T]`: gets a copy of the front of a list without
      removing that range.

- [ ] `Back(count int) List[T]`: gets a copy of the back of a list without
      removing that range.

## Improvements

- [ ] Add atomics and locking to make these data structures optionally
      thread safe.

- [ ] Sorted dictionaries are a very simple implementation which could be
      greatly improved upon. The current solution was just to get things up and
      running. Doing injection sort for keys is not the best solution.
      Instead this should use something like a Red-Black tree to keep the keys
      sorted.

- [ ] Add shortcuts to iterator. Each iterator function should check if the
      iterator implements a method to perform the function faster. For example
      list could create a custom iterator which implements `Next` and `Current`
      but also implements `Count` which returns the count without iterating.
      Then when the iterator `Count` method is called, it checks for the
      shortcut before doing the default. This could improve performance of
      `Reverse`, `Sort`, `ToSlice`, etc. Some shortcuts could also be used
      for other functions, like any algorithm which is improved by knowing
      the count ahead of time.

- [ ] Update the `ToSlice` method based on benchmarks around linked chunks.

- [ ] Improve the `utils.Equal` by doing full recursion as much as possible to
      handle `Equatable` deeper.

- [ ] Add `TError`s for panics in dictionaries and other object for methods
      that can have runtime panics based on inputs. Check inputs before and
      panic mre specific TErrors instead.

- [ ] Add more benchmarks and prove computational complexity.

- [ ] Add better and more documentation.

- [ ] Add more examples.
