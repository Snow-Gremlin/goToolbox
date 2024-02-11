package diff

import (
	"strconv"
	"strings"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/differs"
	"github.com/Snow-Gremlin/goToolbox/differs/data"
)

func runBenchmarks(b *testing.B, data differs.Data, suffix string) {
	b.Run(`Hirschberg-NoReduce`+suffix,
		func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				Hirschberg(-1, false).Diff(data)
			}
		})

	b.Run(`Hirschberg-UseReduce`+suffix,
		func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				Hirschberg(-1, true).Diff(data)
			}
		})

	b.Run(`Wagner`+suffix,
		func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				Wagner(-1).Diff(data)
			}
		})

	b.Run(`Hybrid-NoReduce-100`+suffix,
		func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				Hybrid(-1, false, 100).Diff(data)
			}
		})

	b.Run(`Hybrid-UesReduce-100`+suffix,
		func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				Hybrid(-1, true, 100).Diff(data)
			}
		})

	b.Run(`Hybrid-NoReduce-300`+suffix,
		func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				Hybrid(-1, false, 300).Diff(data)
			}
		})

	b.Run(`Hybrid-UseReduce-300`+suffix,
		func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				Hybrid(-1, true, 300).Diff(data)
			}
		})

	b.Run(`Hybrid-NoReduce-500`+suffix,
		func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				Hybrid(-1, false, 500).Diff(data)
			}
		})

	b.Run(`Hybrid-UseReduce-500`+suffix,
		func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				Hybrid(-1, true, 500).Diff(data)
			}
		})
}

func Benchmark_Simple_Comparison(b *testing.B) {
	comp := data.Strings(exampleA, exampleB)
	runBenchmarks(b, comp, ``)
}

func Benchmark_Default_Reuse(b *testing.B) {
	comp := data.Strings(strings.Split(billNyeA, ` `), strings.Split(billNyeB, ` `))

	b.Run(`Default-NoReused`, func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			Default().Diff(comp)
		}
	})

	b.Run(`Default-Reused`, func(b *testing.B) {
		diff := Default()
		for n := 0; n < b.N; n++ {
			diff.Diff(comp)
		}
	})
}

func Benchmark_Basic_Comparison(b *testing.B) {
	const groups = 4
	for i := 0; i < groups; i++ {
		inputA := billNyeA[:len(billNyeA)*(i+1)/groups]
		inputB := billNyeB[:len(billNyeB)*(i+1)/groups]
		comp := data.Chars(inputA, inputB)
		suffix := `-` + strconv.Itoa(len(inputA)*len(inputB))
		runBenchmarks(b, comp, suffix)
	}
}

func Benchmark_Variant_Comparison(b *testing.B) {
	const groups = 4
	for i := 0; i < groups; i++ {
		inputA := billNyeA[:len(billNyeA)*(i+1)/groups]
		for j := 0; j < groups; j++ {
			inputB := billNyeB[:len(billNyeB)*(j+1)/groups]
			comp := data.Chars(inputA, inputB)
			suffix := `-` + strconv.Itoa(len(inputA)*len(inputB))
			runBenchmarks(b, comp, suffix)
		}
	}
}
