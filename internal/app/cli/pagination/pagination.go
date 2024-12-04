package pagination

import (
	"fmt"

	"workshop-1/internal/domain"
)

func PrintReturns(returns []domain.OrderReturn) error {
	if len(returns) == 0 {
		return fmt.Errorf("не было получено ни одного возврата")
	}

	fmt.Println("=====Возвраты пользователя=====")
	for _, oReturn := range returns {
		fmt.Printf("%+v\n", oReturn)
	}

	return nil
}
