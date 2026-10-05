package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func checkFastMul() error {
	random := rand.New(rand.NewSource(42))

	// Проверяем характерные случаи и случайные множители
	multipliers := []int{1, 2, 12, 15, 16, 100}
	for i := 0; i < 14; i++ {
		multipliers = append(multipliers, random.Intn(10_000)+1)
	}

	var program strings.Builder
	program.WriteString("package main\n\nimport \"fmt\"\n\n")

	for i, y := range multipliers {
		source := fastMulSource(y)

		// Переименовываем функции, чтобы собрать их в одну программу
		source = strings.Replace(
			source,
			"func f(",
			fmt.Sprintf("func f%d(", i),
			1,
		)

		program.WriteString(source)
		program.WriteString("\n")
	}

	program.WriteString("func main() {\n")

	for i, y := range multipliers {
		// 20 множителей × 500 входов = 10000 проверок
		for j := 0; j < 500; j++ {
			x := random.Intn(20_001) - 10_000

			// включаем ноль и числа обоих знаков
			switch j {
			case 0:
				x = 0
			case 1:
				x = 1
			case 2:
				x = -1
			}

			want := x * y

			fmt.Fprintf(
				&program,
				"\tif got := f%d(%d); got != %d {\n"+
					"\t\tpanic(fmt.Sprintf(%q, got))\n"+
					"\t}\n",
				i,
				x,
				want,
				fmt.Sprintf(
					"x=%d, y=%d: получили %%d, ожидали %d",
					x, y, want,
				),
			)
		}
	}

	program.WriteString("}\n")

	directory, err := os.MkdirTemp("", "fast-mul-check-*")
	if err != nil {
		return fmt.Errorf("создание временной папки: %w", err)
	}
	defer os.RemoveAll(directory)

	path := filepath.Join(directory, "generated.go")
	if err := os.WriteFile(path, []byte(program.String()), 0600); err != nil {
		return fmt.Errorf("запись сгенерированной программы: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	output, err := exec.CommandContext(
		ctx,
		"go",
		"run",
		path,
	).CombinedOutput()

	if ctx.Err() != nil {
		return fmt.Errorf(
			"проверка превысила лимит 60 секунд: %w",
			ctx.Err(),
		)
	}

	if err != nil {
		return fmt.Errorf(
			"проверка сгенерированного кода: %w\n%s",
			err,
			output,
		)
	}

	return nil
}
