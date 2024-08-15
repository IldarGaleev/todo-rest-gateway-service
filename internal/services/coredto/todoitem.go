// Package coredto contains Core DTO models
package coredto

type ToDoItem struct {
	ItemID *uint64
	Title  *string
	IsDone *bool
	Owner  *User
}
