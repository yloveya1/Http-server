package service

import (
	"context"
	"fmt"
	"go.mod/internal/domain"
	"go.mod/pkg/repository"
)

// var er repository.Repository = &service{}
const (
	ordersTable   = "orders"
	deliveryTable = "delivery"
	paymentTable  = "payment"
	itemsTable    = "items"
)

func NewService(client repository.Client) repository.Repository {
	return &Service{client: client}
}

type Service struct {
	client repository.Client
	// logger
}

func (r *Service) FindDelivery(ctx context.Context, delivery *domain.Delivery, uid string) error {
	q := fmt.Sprintf("SELECT name, phone, zip, city, address, region, email FROM %s WHERE order_uid=$1", deliveryTable)
	err := r.client.QueryRow(ctx, q, uid).Scan(delivery.Name, delivery.Phone, delivery.Zip, delivery.City, delivery.Address, delivery.Region, delivery.Email)
	if err != nil {
		//
	}
	return err
}

func (r *Service) FindItems(ctx context.Context, items *[]domain.Item, uid string) error {
	q := fmt.Sprintf("SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM %s WHERE order_uid=$1", itemsTable)

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return err
	}

	for rows.Next() {
		var item domain.Item

		err = rows.Scan(&item.Chrt_id, &item.Track_number, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.Total_price, &item.Nm_id, &item.Brand, &item.Status)
		if err != nil {
			return nil
		}

		*items = append(*items, item)
	}
	return nil
}

func (r *Service) FindPayment(ctx context.Context, payment *domain.Payment, uid string) error {
	q := fmt.Sprintf("SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM %s WHERE order_uid=$1", paymentTable)
	err := r.client.QueryRow(ctx, q, uid).Scan(&payment.Transaction, &payment.Request_id, &payment.Currency, &payment.Provider, &payment.Amount, &payment.Payment_dt, &payment.Bank, &payment.Delivery_cost, &payment.Goods_total, &payment.Custom_fee)
	if err != nil {
		//
	}
	return err
}

func (r *Service) Create(ctx context.Context, order *domain.Order) error {
	var err error
	return err
}

func (r *Service) FindOrder(ctx context.Context, uid string) (domain.Order, error) {
	//TODO implement me
	var ord domain.Order

	q := fmt.Sprintf("SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM %s WHERE order_uid=$1", ordersTable)

	err := r.client.QueryRow(ctx, q, uid).Scan(&ord.Order_uid, &ord.Track_number, &ord.Entry, &ord.Locale, &ord.Internal_signature, &ord.Customer_id, &ord.Delivery_service, &ord.Shardkey, &ord.Sm_id, &ord.Date_created, &ord.Oof_shard)
	if err != nil {
		//
	}
	err = r.FindPayment(ctx, &ord.Payment, uid)
	if err != nil {
		//
	}
	err = r.FindDelivery(ctx, &ord.Delivery, uid)
	if err != nil {
		//
	}
	err = r.FindItems(ctx, &ord.Items, uid)
	if err != nil {
		//
	}
	return ord, nil
}

func (r *Service) FindAllData(ctx context.Context) (map[string]domain.Order, error) {
	q := fmt.Sprintf("SELECT order_uid FROM %s", ordersTable)
	var orders map[string]domain.Order
	rows, err := r.client.Query(ctx, q)

	for rows.Next() {
		var temp string
		err = rows.Scan(&temp)
		if err != nil {
			return nil, err
		}
		orders[temp], err = r.FindOrder(ctx, temp)
		if err != nil {
			return nil, err
		}
	}
	return orders, nil
}
