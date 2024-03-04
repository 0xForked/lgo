package todo

type TODO struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

var TODOs []TODO

func (todo TODO) Fetch() []TODO {
	return TODOs
}

func (todo TODO) Create(body string) TODO {
	todo.ID = len(TODOs) + 1
	todo.Body = body
	TODOs = append(TODOs, todo)
	return todo
}
