package main

import (
	"database/sql"
	//"time"
)

type PackageRepository struct {
	db *sql.DB
}

func NewPackageRepository(db *sql.DB) *PackageRepository {
	return &PackageRepository{db: db}
}

// ======================
// CREATE PACKAGE (FIXED)
// ======================
func (r *PackageRepository) Create(pkg *Package) error {

	_, err := r.db.Exec(
		`INSERT INTO packages
		(resi, nama_barang, berat, warehouse_zone, status)
		VALUES (?, ?, ?, ?, ?)`,
		pkg.Resi,
		pkg.NamaBarang,
		pkg.Berat,
		pkg.WarehouseZone,
		pkg.Status,
	)

	return err
}

// ======================
// GET BY RESI (FIXED)
// ======================
func (r *PackageRepository) GetByResi(resi string) (*Package, error) {

	var pkg Package

	err := r.db.QueryRow(
		`SELECT
			resi,
			nama_barang,
			berat,
			warehouse_zone,
			status,
			sorted_at
		FROM packages
		WHERE resi = ?`,
		resi,
	).Scan(
		&pkg.Resi,
		&pkg.NamaBarang,
		&pkg.Berat,
		&pkg.WarehouseZone,
		&pkg.Status,
		&pkg.SortedAt,
	)

	if err != nil {
		return nil, err
	}

	return &pkg, nil
}

// COMPLETE SORT
func (r *PackageRepository) CompleteSort(resi string) error {

	_, err := r.db.Exec(
		`UPDATE packages
		SET status='ready',
		    sorted_at=NOW()
		WHERE resi=?`,
		resi,
	)

	return err
}

// OUTBOX
func (r *PackageRepository) SaveOutbox(eventType string, payload string) error {
	_, err := r.db.Exec(
		`INSERT INTO outbox_events (event_type, payload, status)
		 VALUES (?, ?, 'pending')`,
		eventType,
		payload,
	)
	return err
}

// UPDATE STATUS
func (r *PackageRepository) UpdateStatus(resi string, status string) error {
	_, err := r.db.Exec(
		`UPDATE packages SET status=? WHERE resi=?`,
		status,
		resi,
	)
	return err
}

// OUTBOX FETCH
func (r *PackageRepository) GetPendingEvents() ([]OutboxEvent, error) {

	rows, err := r.db.Query(
		`SELECT id, event_type, payload
		 FROM outbox_events
		 WHERE status='pending'`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []OutboxEvent

	for rows.Next() {
		var e OutboxEvent
		rows.Scan(&e.ID, &e.EventType, &e.Payload)
		events = append(events, e)
	}

	return events, nil
}

// MARK SENT
func (r *PackageRepository) MarkAsSent(id int) error {
	_, err := r.db.Exec(
		`UPDATE outbox_events
		 SET status='sent', sent_at=NOW()
		 WHERE id=?`,
		id,
	)
	return err
}

func (r *PackageRepository) GetAll() ([]Package, error) {

	rows, err := r.db.Query(`
		SELECT
			resi,
			nama_barang,
			berat,
			warehouse_zone,
			status,
			sorted_at
		FROM packages
		ORDER BY resi
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var packages []Package

	for rows.Next() {
		var p Package

		err := rows.Scan(
			&p.Resi,
			&p.NamaBarang,
			&p.Berat,
			&p.WarehouseZone,
			&p.Status,
			&p.SortedAt,
		)

		if err != nil {
			return nil, err
		}

		packages = append(packages, p)
	}

	return packages, nil
}