package main

import (
	"fmt"
	"math/bits"
	"strings"
)

func fastMul(x, y int) int {
	r := x

	for bit := bits.Len(uint(y)) - 2; bit >= 0; bit-- {
		r = r + r

		if uint(y)&(uint(1)<<bit) != 0 {
			r = r + x
		}
	}

	return r
}

func fastMulSource(y int) string {
	var source strings.Builder

	source.WriteString("func f(x int) int {\n")
	source.WriteString("\tr := x\n")

	for bit := bits.Len(uint(y)) - 2; bit >= 0; bit-- {
		// Каждый следующий бит требует удвоения результата
		source.WriteString("\tr = r + r\n")

		// Только единичный бит требует дополнительного прибавления x
		if uint(y)&(uint(1)<<bit) != 0 {
			source.WriteString("\tr = r + x\n")
		}
	}

	source.WriteString("\treturn r\n")
	source.WriteString("}\n")

	return source.String()
}


func fastMulGen(y int) {
	fmt.Println(fastMulSource(y))
}

func main() {
	var y int

	fmt.Print("Введите положительный множитель y: ")
	if _, err := fmt.Scan(&y); err != nil || y < 1 {
		fmt.Println("Ошибка: нужно ввести целое число больше нуля")
		return
	}

	fmt.Println("\nСгенерированная функция:")
	fastMulGen(y)

	if err := checkFastMul(); err != nil {
		fmt.Println("\nАвтоматическая проверка не прошла:", err)
		return
	}

	fmt.Println("\nАвтоматическая проверка: 10000 примеров пройдено")
}
