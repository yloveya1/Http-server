package service

import (
	"context"
	"fmt"
	"go.mod/internal/cache"
	"go.mod/internal/domain"
	"go.mod/pkg/repository"
	"time"
)

const (
	ordersTable   = "orders"
	deliveryTable = "delivery"
	paymentTable  = "payment"
	itemsTable    = "items"
)

func NewService(client repository.Client, cache *cache.Cache) repository.Repository {
	return &Service{client: client, cache: cache}
}

type Service struct {
	client repository.Client
	cache  *cache.Cache
}

func (r *Service) AddOrderDataDB(ctx context.Context, order domain.Order) error {
	//Insert order
	q := fmt.Sprintf(`INSERT INTO %s(
			order_uid, track_number, entry, locale, internal_signature, 
            customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING order_uid;`, ordersTable)

	if err := r.client.QueryRow(ctx, q, order.Order_uid, order.Track_number, order.Entry, order.Locale, order.Internal_signature, order.Customer_id, order.Delivery_service, order.Shardkey, order.Sm_id, order.Date_created, order.Oof_shard).Scan(&order.Order_uid); err != nil {
		return err
	}

	// Insert delivery
	q = fmt.Sprintf(`INSERT INTO %s(
			order_uid, name, phone, zip, city, address, region, email)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING order_uid;`, deliveryTable)

	if err := r.client.QueryRow(ctx, q, order.Order_uid, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email).Scan(&order.Order_uid); err != nil {
		return err
	}
	//Insert payment
	q = fmt.Sprintf(`INSERT INTO %s(
			order_uid, transaction, request_id, currency, provider, amount, 
            payment_dt, bank, delivery_cost, goods_total, custom_fee)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING order_uid;`, paymentTable)
	if err := r.client.QueryRow(ctx, q, order.Order_uid, order.Payment.Transaction, order.Payment.Request_id,
		order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.Payment_dt,
		order.Payment.Bank, order.Payment.Delivery_cost, order.Payment.Goods_total, order.Payment.Custom_fee).Scan(&order.Order_uid); err != nil {
		return err
	}
	// Insert items

	q = fmt.Sprintf(`INSERT INTO %s(
			order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING order_uid;`, itemsTable)
	for _, val := range order.Items {
		if err := r.client.QueryRow(ctx, q, order.Order_uid, val.Chrt_id, val.Track_number, val.Price, val.Rid, val.Name,
			val.Sale, val.Size, val.Total_price, val.Nm_id, val.Brand, val.Status).Scan(&order.Order_uid); err != nil {
			return err
		}
	}
	return nil
}

func (r *Service) FindDelivery(ctx context.Context, delivery *domain.Delivery, uid string) error {
	q := fmt.Sprintf("SELECT name, phone, zip, city, address, region, email FROM %s WHERE order_uid=$1", deliveryTable)
	if err := r.client.QueryRow(ctx, q, uid).Scan(&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address, &delivery.Region, &delivery.Email); err != nil {
		return err
	}
	return nil
}

func (r *Service) FindItems(ctx context.Context, items *domain.Items, uid string) error {
	q := fmt.Sprintf("SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM %s WHERE order_uid=$1", itemsTable)

	rows, err := r.client.Query(ctx, q, uid)
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
	if err := r.client.QueryRow(ctx, q, uid).Scan(&payment.Transaction, &payment.Request_id, &payment.Currency, &payment.Provider, &payment.Amount, &payment.Payment_dt, &payment.Bank,
		&payment.Delivery_cost, &payment.Goods_total, &payment.Custom_fee); err != nil {
		return err
	}
	return nil
}

func (r *Service) FindOrder(ctx context.Context, uid string) (domain.Order, error) {
	//TODO implement me
	var ord domain.Order

	q := fmt.Sprintf("SELECT order_uid, track_number, entry, locale, internal_signature,"+
		" customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM %s WHERE order_uid=$1", ordersTable)

	err := r.client.QueryRow(ctx, q, uid).Scan(&ord.Order_uid, &ord.Track_number, &ord.Entry, &ord.Locale,
		&ord.Internal_signature, &ord.Customer_id, &ord.Delivery_service, &ord.Shardkey, &ord.Sm_id, &ord.Date_created,
		&ord.Oof_shard)

	if err != nil {
		return domain.Order{}, err
	}
	err = r.FindPayment(ctx, &ord.Payment, uid)
	if err != nil {
		return domain.Order{}, err
	}
	err = r.FindDelivery(ctx, &ord.Delivery, uid)
	if err != nil {
		return domain.Order{}, err
	}
	err = r.FindItems(ctx, &ord.Items, uid)
	if err != nil {
		return domain.Order{}, err
	}
	return ord, nil
}

func (r *Service) FindAllOrders(ctx context.Context) error {
	q := fmt.Sprintf("SELECT order_uid FROM %s", ordersTable)
	orders := map[string]domain.Order{}
	rows, err := r.client.Query(ctx, q)
	for rows.Next() {
		var temp string
		err = rows.Scan(&temp)
		if err != nil {
			return err
		}
		orders[temp], err = r.FindOrder(ctx, temp)
		r.cache.Set(temp, orders[temp], 10*time.Minute)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Service) GetCache() *cache.Cache {
	return r.cache
}
