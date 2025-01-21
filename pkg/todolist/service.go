package todolist

import (
	"context"

	"go.altair.com/todolist/pkg/structs"
	"go.altair.com/todolist/pkg/todolist/store"
)

type ItemsService interface {
	AddItem(ctx context.Context, def *structs.TodoItem) error
	DeleteItem(ctx context.Context, id string) error
	UpdateItem(ctx context.Context, def *structs.TodoItem) error
	GetItem(ctx context.Context, id string) (*structs.TodoItem, error)
	ListItems(ctx context.Context) (structs.TodoItemList, error)
}

func NewItemsService(s store.Store) ItemsService {
	return &itemsServiceImpl{
		store: s,
	}
}

type itemsServiceImpl struct {
	store store.Store
}

func (s *itemsServiceImpl) GetItem(ctx context.Context, deploymentId string) (*structs.TodoItem, error) {
	var result structs.TodoItem
	err := s.store.Update(func(tx store.Txn) error {
		err := tx.Get(ctx, deploymentId, &result)
		return err
	})
	return &result, err
}

func (s *itemsServiceImpl) AddItem(ctx context.Context, def *structs.TodoItem) error {
	return s.store.Update(func(tx store.Txn) error {
		
		return tx.Add(ctx, def)
	})
}

func (s *itemsServiceImpl) ListItems(ctx context.Context) (structs.TodoItemList, error) {
	var result structs.TodoItemList
	err := s.store.Update(func(tx store.Txn) error {
		err := tx.List(ctx, &result)
		
		return err
	})
	return result, err
}

func (s *itemsServiceImpl) DeleteItem(ctx context.Context, deploymentId string) error {
	return s.store.Update(func(tx store.Txn) error {
		return tx.Delete(ctx, deploymentId)
	})
}

func (s *itemsServiceImpl) UpdateItem(ctx context.Context, def *structs.TodoItem) error {
	return s.store.Update(func(tx store.Txn) error {
		return tx.Update(ctx, def)
	})
}

