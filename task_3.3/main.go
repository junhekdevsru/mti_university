package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func highestBitMask(y int) int {
	mask := 1

	for mask <= y/2 {
		mask <<= 1
	}

	return mask
}

// fastMul — обычное двоичное умножение.
// Начинаем с коэффициента 1, то есть result = x.
func fastMul(x, y int) int {
	if y < 1 {
		panic("multiplier must be positive")
	}

	result := x
	mask := highestBitMask(y) >> 1

	for mask > 0 {
		result = result + result

		if y&mask != 0 {
			result = result + x
		}

		mask >>= 1
	}

	return result
}

func isPowerOfTwo(value int) bool {
	return value > 0 && value&(value-1) == 0
}

// isMersenne определяет числа вида 2^n - 1.
//
// Например:
//
//	7  = 2^3 - 1
//	15 = 2^4 - 1
//	31 = 2^5 - 1
func isMersenne(y int) bool {
	return y > 0 && isPowerOfTwo(y+1)
}

// writeBinaryFormula генерирует обычную двоичную формулу.
func writeBinaryFormula(source *strings.Builder, y int) {
	source.WriteString("\tr := x\n")

	mask := highestBitMask(y) >> 1

	for mask > 0 {
		source.WriteString("\tr = r + r\n")

		if y&mask != 0 {
			source.WriteString("\tr = r + x\n")
		}

		mask >>= 1
	}
}

// writeMersenneFormula использует формулу для чисел вида 2^n - 1.
//
// Для y = 15:
//
//	15 = 8 - (1 - 8)
//
// Если r = 8x, то:
//
//	negative = x - r = -7x
//	r = r - negative = 8x - (-7x) = 15x
//
// Получается 3 сложения и 2 вычитания.
func writeMersenneFormula(source *strings.Builder, y int) {
	source.WriteString("\tr := x\n")

	half := (y + 1) / 2
	coefficient := 1

	for coefficient < half {
		source.WriteString("\tr = r + r\n")
		coefficient <<= 1
	}

	source.WriteString("\tnegative := x - r\n")
	source.WriteString("\tr = r - negative\n")
}

func fastMulSource(y int) (string, error) {
	if y < 1 {
		return "", fmt.Errorf("multiplier must be positive")
	}

	var source strings.Builder

	source.WriteString("func f(x int) int {\n")

	// Для 15, 31, 63 и других достаточно больших чисел вида 2^n - 1
	// используем более короткую формулу со сложениями и вычитаниями.
	if isMersenne(y) && y >= 15 {
		writeMersenneFormula(&source, y)
	} else {
		writeBinaryFormula(&source, y)
	}

	source.WriteString("\treturn r\n")
	source.WriteString("}\n")

	return source.String(), nil
}

// fastMulGen непосредственно печатает сгенерированную функцию.
func fastMulGen(y int) error {
	source, err := fastMulSource(y)
	if err != nil {
		return err
	}

	fmt.Print(source)
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(
			os.Stderr,
			"usage: %s <positive multiplier>\n",
			os.Args[0],
		)
		os.Exit(2)
	}

	y, err := strconv.Atoi(os.Args[1])
	if err != nil || y < 1 {
		fmt.Fprintln(os.Stderr, "multiplier must be a positive integer")
		os.Exit(2)
	}

	if err := fastMulGen(y); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
