package tests

const (
	TODO_CREATE_QUERY         string = `INSERT INTO todos (id,user_id,content,status) VALUES ($1,$2,$3,$4) RETURNING id, created_at`
	TODO_GET_TODOS_BY_USER_ID string = `SELECT id, user_id, content, status, created_at, updated_at FROM todos`
	TODO_GET_TODO_BY_USER_ID  string = `SELECT id, user_id, content, status, created_at, updated_at FROM todos WHERE user_id = $1 AND id = $2`
	TODO_UPDATE_STATUS        string = `UPDATE todos SET status = $1 WHERE id = $2 AND user_id = $3`
	TODO_UPDATE_CONTENT       string = `UPDATE todos SET content = $1 WHERE id = $2 AND user_id = $3`
	TODO_DELETE               string = `DELETE FROM todos WHERE id = $1 AND user_id = $2`

	USER_GET_BY_EMAIL string = `SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1`
)
