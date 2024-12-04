package postgres

type Storage struct {
	txManager TransactionManager
}

func NewStorage(txManager TransactionManager) *Storage {
	return &Storage{
		txManager: txManager,
	}
}
