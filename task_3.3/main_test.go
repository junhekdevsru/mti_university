package main

import (
	"go/parser"
	"go/token"
	"math/rand"
	"strings"
	"testing"
)

func TestFastMulRandomValues(t *testing.T) {
	random := rand.New(rand.NewSource(42))

	for i := 0; i < 10_000; i++ {
		x := random.Intn(2_000_001) - 1_000_000
		y := random.Intn(100_000)

		got := fastMul(x, y)
		want := x * y

		if got != want {
			t.Fatalf(
				"fastMul(%d, %d) = %d; want %d",
				x,
				y,
				got,
				want,
			)
		}
	}
}

func TestFastMulSpecialCases(t *testing.T) {
	tests := []struct {
		x    int
		y    int
		want int
	}{
		{x: 10, y: 0, want: 0},
		{x: 10, y: 1, want: 10},
		{x: 10, y: 12, want: 120},
		{x: 10, y: 15, want: 150},
		{x: 10, y: 16, want: 160},
		{x: -10, y: 15, want: -150},
	}

	for _, test := range tests {
		got := fastMul(test.x, test.y)

		if got != test.want {
			t.Errorf(
				"fastMul(%d, %d) = %d; want %d",
				test.x,
				test.y,
				got,
				test.want,
			)
		}
	}
}

func TestGeneratedProgramsHaveValidSyntax(t *testing.T) {
	multipliers := []int{0, 1, 12, 15, 16, 255}

	for _, y := range multipliers {
		source := fastMulSource(y)
		fileSource := "package generated\n\n" + source

		_, err := parser.ParseFile(
			token.NewFileSet(),
			"generated.go",
			fileSource,
			parser.AllErrors,
		)
		if err != nil {
			t.Errorf(
				"invalid generated code for y=%d:\n%s\nerror: %v",
				y,
				source,
				err,
			)
		}

		if strings.Contains(source, "*") {
			t.Errorf(
				"generated function for y=%d contains multiplication:\n%s",
				y,
				source,
			)
		}
	}
}

func TestGeneratedAdditionCount(t *testing.T) {
	tests := []struct {
		y    int
		want int
	}{
		{y: 12, want: 5},
		{y: 16, want: 5},
		{y: 15, want: 7},
	}

	for _, test := range tests {
		source := fastMulSource(test.y)
		got := strings.Count(source, " + ")

		if got != test.want {
			t.Errorf(
				"generated function for y=%d uses %d additions; want %d",
				test.y,
				got,
				test.want,
			)
		}
	}
}

func TestFastMulSource15(t *testing.T) {
	got := fastMulSource(15)

	want := `func f(x int) int {
		a := x + x // 2x — первое сложение
		b := a + a // 4x — второе сложение
		c := b + b // 8x — третье сложение
		d := x - c // -7x — первое вычитание
		r := c - d // 15x — второе вычитание
		return r
	}
`

	if got != want {
		t.Fatalf("unexpected generated code:\n%s\nwant:\n%s", got, want)
	}
}
