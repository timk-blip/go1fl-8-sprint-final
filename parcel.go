package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	client := p.Client
	status := p.Status
	address := p.Address
	creatAt := p.CreatedAt
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", client),
		sql.Named("status", status),
		sql.Named("address", address),
		sql.Named("created_at", creatAt))
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	var parcel Parcel
	row, err := s.db.Query("SELECT * FROM parcel WHERE number = $1", number)
	if err != nil {
		return parcel, err
	}
	defer row.Close()
	// заполните объект Parcel данными из таблицы
	p := Parcel{}
	if !row.Next() {
		return p, fmt.Errorf("parcel not found")
	}
	err = row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = $1", client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []Parcel
	for rows.Next() {
		p := Parcel{}
		err = rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	// заполните срез Parcel данными из таблицы
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	p, err := s.Get(number)
	if err != nil {
		return err
	}
	if p.Status == "registered" {
		_, err = s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number",
			sql.Named("address", address),
			sql.Named("number", number))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	p, err := s.Get(number)
	if err != nil {
		return err
	}
	if p.Status == "registered" {
		_, err = s.db.Exec("DELETE FROM parcel WHERE number = :number",
			sql.Named("number", number))
		if err != nil {
			return err
		}
	}
	return nil
}
