.PHONY: run-server get-todos create-todo

run-server:
	go run ./main.go  -port="8000" -env="local"

get-todos:
	curl -X GET http://localhost:8000/todos

# call: make create-todo todo="Buy groceries"
create-todo:
	curl -X POST -d '{"body": "$(todo)"}' http://localhost:8000/todos
