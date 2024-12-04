package pagination

import (
	"fmt"
	"os"
	"syscall"

	"workshop-1/internal/domain"

	"golang.org/x/term"
)

func ScrollOrders(orders []domain.Order, page, pageSize int) error {
	state, err := term.MakeRaw(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("ошибка переключения консоли: %w", err)
	}
	defer term.Restore(int(syscall.Stdin), state)

LOOP:
	for {
		printPage(orders, page, pageSize)

		buf := make([]byte, 3)
		os.Stdin.Read(buf)

		switch string(buf) {
		case "\x1b[B": // стрелка вниз
			page = min(len(orders)/pageSize+1, page+1)
		case "\x1b[A": // стрелка вверх
			page = max(1, page-1)
		default:
			break LOOP
		}
	}

	return nil
}

func printPage(orders []domain.Order, page, pageSize int) {
	// очистка консоли
	fmt.Print("\033[H\033[2J")

	fmt.Printf("\nИспользуйте стрелки вверх и вниз, чтобы листать список")
	fmt.Printf("\nНажмите любую клавишу, чтобы выйти\n\n")

	if page == 1 {
		fmt.Printf("=====Начало списка=====\n")
	}

	for i := (page - 1) * pageSize; i < page*pageSize && i < len(orders); i++ {
		fmt.Printf("%+v\n", orders[i])
	}

	if page*pageSize >= len(orders) {
		fmt.Printf("=====Конец списка=====\n")
	}
}
