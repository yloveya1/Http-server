package repository

import (
	"context"
	"go.mod/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, order *domain.Order) error
	FindOrder(ctx context.Context, uid string) (domain.Order, error)
	FindDelivery(ctx context.Context, delivery *domain.Delivery, uid string) error
	FindItems(ctx context.Context, items *[]domain.Item, uid string) error
	FindPayment(ctx context.Context, payment *domain.Payment, uid string) error
	FindAllData(ctx context.Context) (map[string]domain.Order, error)
}
