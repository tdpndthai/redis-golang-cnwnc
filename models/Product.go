package models

import (
	"database/sql"
	"redis-golang-cnwnc/entities"
)

type ProductModel struct {
	Db *sql.DB
}

func (productModel ProductModel) GetAll() (proc []entities.Product, err error) { //method truyền vào procmodel trả về mảng proc
	rows, err := productModel.Db.Query("select * from proc")
	if err != nil {
		return nil, err
	} else {
		var procs []entities.Product
		for rows.Next() {
			var id int64
			var name string
			var price float64
			var quantity int64
			err2 := rows.Scan(&id, &name, &price, &quantity) //sao chép các cột trong hàng hiện tại vào các giá trị được chỉ định
			if err2 != nil {
				return nil, err2
			} else {
				proc := entities.Product{
					ID:       id,
					Name:     name,
					Price:    price,
					Quantity: quantity,
				}
				procs = append(procs, proc)
			}
		}
		return procs, nil
	}
}
