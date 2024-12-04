package cobra

import (
	"context"
	"fmt"
	"time"

	"workshop-1/internal/app/cli/pagination"
	"workshop-1/internal/domain/strategy"
	"workshop-1/internal/dto"
	"workshop-1/internal/usecase"

	"github.com/spf13/cobra"
)

const (
	orderFlag  = "order"
	ordersFlag = "orders"
	userFlag   = "user"
	expiryFlag = "expiry"
	weightFlag = "weight"
	priceFlag  = "price"

	packagingFlag = "packaging"
	addWrapFlag   = "add-wrap"

	lastFlag     = "last"
	inPVZFlag    = "in-pvz"
	pageFlag     = "page"
	pageSizeFlag = "pageSize"

	goroutinesNumberFlag = "num"
)

func newCourierAcceptCmd(ctx context.Context, storage usecase.Storage) *cobra.Command {
	order := dto.CreateOrder{}
	var packagingStr string
	var packaging strategy.Packaging
	var addWrap bool

	cmd := &cobra.Command{
		Use:   "courier-accept",
		Short: "Accept an order from the courier",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			pkg, err := strategy.NewPackaging(packagingStr)
			if err != nil {
				return err
			}

			if addWrap {
				pkg = strategy.PackagingWithWrap{MainPackaging: pkg}
			}

			packaging = pkg

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			currentTime := time.Now().UTC()
			return usecase.AcceptFromCourier(ctx, storage, currentTime, order, packaging)
		},
	}

	cmd.Flags().IntVar(&order.ID, orderFlag, 0, "order id (required)")
	cmd.Flags().IntVar(&order.UserID, userFlag, 0, "user id (required)")
	cmd.Flags().StringVar(&order.ExpiryDate, expiryFlag, "", "expiry date dd-mm-yyyy (required)")
	cmd.Flags().Float32Var(&order.Weight, weightFlag, 0, "weight (required)")
	cmd.Flags().Float32Var(&order.Price, priceFlag, 0, "price (required)")

	cmd.MarkFlagRequired(orderFlag)
	cmd.MarkFlagRequired(userFlag)
	cmd.MarkFlagRequired(expiryFlag)
	cmd.MarkFlagRequired(weightFlag)
	cmd.MarkFlagRequired(priceFlag)

	cmd.Flags().StringVar(&packagingStr, packagingFlag, "", "packaging: wrap, bag, box (required)")
	cmd.Flags().BoolVar(&addWrap, addWrapFlag, false, "add wrap (optional)")

	cmd.MarkFlagRequired(packagingFlag)

	return cmd
}

func newCourierReturnCmd(ctx context.Context, storage usecase.Storage) *cobra.Command {
	var orderID int

	cmd := &cobra.Command{
		Use:   "courier-return",
		Short: "Return an order to the courier",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if orderID < 0 {
				return fmt.Errorf("неверный ID заказа: %d", orderID)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			currentTime := time.Now().UTC()
			return usecase.ReturnToCourier(ctx, storage, currentTime, orderID)
		},
	}

	cmd.Flags().IntVar(&orderID, orderFlag, 0, "order id (required)")
	cmd.MarkFlagRequired(orderFlag)

	return cmd
}

func newUserGiveCmd(ctx context.Context, storage usecase.Storage) *cobra.Command {
	var orderIDs []int

	cmd := &cobra.Command{
		Use:   "user-give",
		Short: "Give orders to the user",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			for _, orderID := range orderIDs {
				if orderID < 0 {
					return fmt.Errorf("неверный ID заказа: %d", orderID)
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			currentTime := time.Now().UTC()
			return usecase.GiveToUser(ctx, storage, currentTime, orderIDs)
		},
	}

	cmd.Flags().IntSliceVar(&orderIDs, ordersFlag, []int{}, "list of order IDs (required)")
	cmd.MarkFlagRequired(ordersFlag)

	return cmd
}

func newUserOrdersCmd(ctx context.Context, storage usecase.Storage) *cobra.Command {
	var userID, lastN int
	var inPVZOnly bool

	cmd := &cobra.Command{
		Use:   "user-orders",
		Short: "List user orders",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if userID < 0 {
				return fmt.Errorf("неверный ID пользователя: %d", userID)
			}
			if lastN < 0 {
				return fmt.Errorf("неверное число последних заказов: %d", lastN)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			pageSize := 5
			page := 1

			orders, err := usecase.UserOrders(ctx, storage, userID, lastN, inPVZOnly)
			if err != nil {
				return err
			}

			return pagination.ScrollOrders(orders, page, pageSize)
		},
	}

	cmd.Flags().IntVar(&userID, userFlag, 0, "user id (required)")
	cmd.Flags().IntVar(&lastN, lastFlag, 0, "number of last orders to get (optional)")
	cmd.Flags().BoolVar(&inPVZOnly, inPVZFlag, false, "only orders that are still in this pickup point (optional)")

	cmd.MarkFlagRequired(userFlag)

	return cmd
}

func newUserReturnCmd(ctx context.Context, storage usecase.Storage) *cobra.Command {
	var userID, orderID int
	cmd := &cobra.Command{
		Use:   "user-return",
		Short: "Accept the return from the user",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if userID < 0 {
				return fmt.Errorf("неверный ID пользователя: %d", userID)
			}
			if orderID < 0 {
				return fmt.Errorf("неверный ID заказа: %d", orderID)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			currentTime := time.Now().UTC()
			return usecase.AcceptUserReturn(ctx, storage, currentTime, userID, orderID)
		},
	}

	cmd.Flags().IntVar(&userID, userFlag, 0, "user id (required)")
	cmd.Flags().IntVar(&orderID, orderFlag, 0, "order id (required)")
	cmd.MarkFlagRequired(userFlag)
	cmd.MarkFlagRequired(orderFlag)

	return cmd
}

func newUserReturnsCmd(ctx context.Context, storage usecase.Storage) *cobra.Command {
	var pageNum, pageSizeInt int
	cmd := &cobra.Command{
		Use:   "user-returns",
		Short: "List user returns",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if pageSizeInt <= 0 {
				return fmt.Errorf("неверный размер страницы: %d", pageSizeInt)
			}
			if pageNum <= 0 {
				return fmt.Errorf("неверный номер страницы: %d", pageNum)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			returns, err := usecase.UserReturns(ctx, storage, pageNum, pageSizeInt)
			if err != nil {
				return err
			}

			return pagination.PrintReturns(returns)
		},
	}

	cmd.Flags().IntVar(&pageNum, pageFlag, 0, "page number (required)")
	cmd.Flags().IntVar(&pageSizeInt, pageSizeFlag, 0, "page size (required)")
	cmd.MarkFlagRequired(pageFlag)
	cmd.MarkFlagRequired(pageSizeFlag)

	return cmd
}

func newChangeGoroutinesCmd(ctx context.Context, c *CobraCLI) *cobra.Command {
	var goroutines int

	cmd := &cobra.Command{
		Use:   "workers",
		Short: "Change number of goroutines",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if goroutines <= 0 {
				return fmt.Errorf("неверное число горутин: %d", goroutines)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c.changeWorkerNum(ctx, goroutines)
			return nil
		},
	}

	cmd.Flags().IntVar(&goroutines, goroutinesNumberFlag, 0, "number of goroutines (required)")
	cmd.MarkFlagRequired(goroutinesNumberFlag)

	return cmd
}
