package repository

import (
	"context"
	"go.mod/internal/cache"
	"go.mod/internal/domain"
)

type Repository interface {
	FindOrder(ctx context.Context, uid string) (domain.Order, error)
	FindDelivery(ctx context.Context, delivery *domain.Delivery, uid string) error
	FindItems(ctx context.Context, items *domain.Items, uid string) error
	FindPayment(ctx context.Context, payment *domain.Payment, uid string) error
	FindAllOrders(ctx context.Context) error
	GetCache() *cache.Cache
	AddOrderDataDB(ctx context.Context, order domain.Order) error
}
