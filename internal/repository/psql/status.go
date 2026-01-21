package psql

type NoteStatus string

const (
	Todo      NoteStatus = "todo"
	Processed NoteStatus = "processed"
	Done      NoteStatus = "done"
)
