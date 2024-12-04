package cobra

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"workshop-1/internal/usecase"

	"github.com/spf13/cobra"
)

type Task struct {
	Args []string
}

type CobraCLI struct {
	rootCmd   *cobra.Command
	workerNum int
	taskChan  chan Task
	mu        sync.Mutex
	wg        sync.WaitGroup
}

func New(ctx context.Context, storage usecase.Storage) *CobraCLI {
	c := &CobraCLI{
		rootCmd:   &cobra.Command{},
		taskChan:  make(chan Task, 2),
		workerNum: 2,
	}

	c.rootCmd.AddCommand(
		newCourierAcceptCmd(ctx, storage),
		newCourierReturnCmd(ctx, storage),
		newUserGiveCmd(ctx, storage),
		newUserOrdersCmd(ctx, storage),
		newUserReturnCmd(ctx, storage),
		newUserReturnsCmd(ctx, storage),
		newChangeGoroutinesCmd(ctx, c),
	)

	c.startWorkers(ctx)
	return c
}

func (c *CobraCLI) Run() error {
	return c.rootCmd.Execute()
}

func (c *CobraCLI) RunArgs(args []string) error {
	c.rootCmd.SetArgs(args)
	return c.rootCmd.Execute()
}

func (c *CobraCLI) RunInteractive(stop context.CancelFunc, reader *bufio.Reader) {
	fmt.Println("Интерактивный режим. Введите 'stop' чтобы выйти.")

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Ошибка при чтении: %v\n", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "stop" {
			fmt.Println("Выход из интерактивного режима...")
			stop()

			return
		}

		args := strings.Fields(input)

		// перехват команды для изменения числа горутин, чтобы не добавлять её как обычную задачу
		if len(args) != 0 && args[0] == "workers" {
			c.RunArgs(args)
			continue
		}

		c.AddTask(Task{Args: args})
	}
}

func (c *CobraCLI) AddTask(task Task) {
	select {
	case c.taskChan <- task:
		fmt.Println("задача добавлена в очередь")
	default:
		fmt.Println("слишком много задач в очереди")
	}
}

func (c *CobraCLI) GracefulShutdown(shutdownCtx context.Context) error {
	fmt.Println("Ожидание завершения задач...")

	close(c.taskChan)

	done := make(chan struct{})
	go func() {
		defer close(done)
		c.wg.Wait()
	}()

	select {
	case <-shutdownCtx.Done():
		return errors.New("вышло время ожидания, досрочное завершение работы...")
	case <-done:
		fmt.Println("Все задачи выполнены, программа завершена.")
		return nil
	}
}
