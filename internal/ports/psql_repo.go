package ports

// создаем задачу
// обновить задачу
// обновить статус согласовальщика
type PSQLPort interface {
	AddEvent()error
}