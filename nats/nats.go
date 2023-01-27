package main

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
	"go.mod/internal/domain"
)

func main() {
	st, err := stan.Connect("test-cluster", "test-client2")
	if err != nil {
		panic(err.Error())
	}
	defer st.Close()
	test := domain.Order{
		Order_uid:          "b563feb7b2b84btest9",
		Customer_id:        "Mytest4",
		Date_created:       "2022-07-24 17:42:39.5555+03:00",
		Delivery_service:   "test4",
		Entry:              "Test4",
		Internal_signature: "24",
		Locale:             "ru",
		Oof_shard:          "12354",
		Shardkey:           "114",
		Sm_id:              24,
		Track_number:       "2554",
		Delivery: domain.Delivery{
			Name:    "Kek4",
			Phone:   "8888888848",
			Zip:     "40808004",
			City:    "Jerusalem4",
			Address: "04",
			Region:  "Deus Vult!4",
			Email:   "god@god.com4"},
		Items: domain.Items{
			domain.Item{
				Chrt_id:      1111,
				Track_number: "1111",
				Price:        1111,
				Rid:          "1111",
				Name:         "1111",
				Sale:         1111,
				Size:         "1111",
				Total_price:  1111,
				Nm_id:        1111,
				Brand:        "1111",
				Status:       1111},
			domain.Item{
				Chrt_id:      2222,
				Track_number: "2222",
				Price:        2222,
				Rid:          "2222",
				Name:         "2222",
				Sale:         2222,
				Size:         "2222",
				Total_price:  2222,
				Nm_id:        2222,
				Brand:        "2222",
				Status:       2222}},
		Payment: domain.Payment{
			Transaction:   "b563feb7b2b84btest4",
			Request_id:    "1234",
			Currency:      "Volt4",
			Provider:      "Sky4",
			Amount:        1234,
			Payment_dt:    1478523694,
			Bank:          "SkyBank4",
			Delivery_cost: 4,
			Goods_total:   999999994,
			Custom_fee:    4}}
	var byts []byte
	byts, err = json.Marshal(&test)
	if err != nil {
		panic(err.Error())
	}
	if err := st.Publish("NewOrder", byts); err != nil {
		panic(err.Error())
	}
}
