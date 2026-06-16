package main

import "database/sql"

type DeliveryRepository struct {
	db *sql.DB
}

func NewDeliveryRepository(
	db *sql.DB,
) *DeliveryRepository {

	return &DeliveryRepository{
		db: db,
	}
}

func (r *DeliveryRepository) Create(
	d *Delivery,
) error {

	_, err := r.db.Exec(
		`INSERT INTO deliveries
		(resi,courier_id,assigned_zone,status)
		VALUES(?,?,?,?)`,
		d.Resi,
		d.CourierID,
		d.AssignedZone,
		d.Status,
	)

	return err
}

func (r *DeliveryRepository) GetAll() ([]Delivery, error) {

	rows, err := r.db.Query(`
		SELECT
			resi,
			courier_id,
			assigned_zone,
			status,
			created_at
		FROM deliveries
		ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deliveries []Delivery

	for rows.Next() {

		var d Delivery

		err := rows.Scan(
			&d.Resi,
			&d.CourierID,
			&d.AssignedZone,
			&d.Status,
			&d.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		deliveries = append(deliveries, d)
	}

	return deliveries, nil
}

func (r *DeliveryRepository) GetByResi(resi string) (*Delivery, error) {

	var d Delivery

	err := r.db.QueryRow(`
		SELECT
			resi,
			courier_id,
			assigned_zone,
			status,
			created_at
		FROM deliveries
		WHERE resi = ?
	`, resi).Scan(
		&d.Resi,
		&d.CourierID,
		&d.AssignedZone,
		&d.Status,
		&d.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &d, nil
}