// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type FetchPost struct {
	ID string `json:"id"`
}

type Mutation struct {
}

type NewPost struct {
	Title string `json:"title"`
}

type Post struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type Query struct {
}
