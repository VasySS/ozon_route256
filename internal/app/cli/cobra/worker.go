package cobra

import (
	"context"
	"fmt"
)

func (c *CobraCLI) changeWorkerNum(ctx context.Context, num int) {
	close(c.taskChan)

	done := make(chan struct{})
	go func() {
		defer close(done)
		c.wg.Wait()
	}()

	fmt.Println("Ожидание завершения задач для изменения количества воркеров...")
	<-done

	c.mu.Lock()
	defer c.mu.Unlock()

	c.workerNum = num
	c.taskChan = make(chan Task, num)
	c.startWorkers(ctx)

	fmt.Printf("Количество воркеров изменено на %d\n", c.workerNum)
}

func (c *CobraCLI) startWorkers(ctx context.Context) {
	for i := 0; i < c.workerNum; i++ {
		go c.worker(ctx)
	}
}

func (c *CobraCLI) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-c.taskChan:
			if !ok {
				return
			}

			c.startTask(task)
		}
	}
}

func (c *CobraCLI) startTask(task Task) {
	c.wg.Add(1)
	defer c.wg.Done()

	c.mu.Lock()
	c.rootCmd.SetArgs(task.Args)
	c.mu.Unlock()

	if err := c.Run(); err != nil {
		fmt.Println(err)
	}

	fmt.Print("> ")
}
