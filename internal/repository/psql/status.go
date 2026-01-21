package psql

type TodoStatus string

const (
	Todo      TodoStatus = "todo"
	Processed TodoStatus = "processed"
	Done      TodoStatus = "done"
)
