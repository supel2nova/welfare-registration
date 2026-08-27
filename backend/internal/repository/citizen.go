package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var ErrCitizenNotFound = errors.New("repository: ไม่พบประชาชนตาม hash")

type AddressParams struct {
	HouseNo         string
	Moo             *string
	Road            *string
	SubdistrictCode string
	DistrictCode    string
	ProvinceCode    string
	PostalCode      string
	SubdistrictName string
	DistrictName    string
	ProvinceName    string
}

type CitizenParams struct {
	NationalIDHash string
	NationalIDEnc  []byte
	Title          string
	FirstName      string
	LastName       string
	BirthYear      int
	BirthMonth     *int
	BirthDay       *int
	BirthPrecision string
	Phone          string
	AddressID      uuid.UUID
}

func (r *Repo) InsertAddress(ctx context.Context, tx pgx.Tx, p AddressParams) (uuid.UUID, error) {
	const q = `
INSERT INTO addresses (house_no, moo, road, subdistrict_code, district_code, province_code,
                       postal_code, subdistrict_name, district_name, province_name)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
RETURNING id`

	var id uuid.UUID
	err := tx.QueryRow(ctx, q, p.HouseNo, p.Moo, p.Road, p.SubdistrictCode, p.DistrictCode,
		p.ProvinceCode, p.PostalCode, p.SubdistrictName, p.DistrictName, p.ProvinceName).Scan(&id)
	return id, err
}

func (r *Repo) UpsertCitizen(ctx context.Context, tx pgx.Tx, p CitizenParams) (uuid.UUID, error) {
	const insert = `
INSERT INTO citizens (national_id_hash, national_id_enc, title, first_name, last_name,
                      birth_year, birth_month, birth_day, birth_precision, phone, address_id)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
ON CONFLICT (national_id_hash) DO NOTHING
RETURNING id`

	const update = `
UPDATE citizens
SET title = $2, first_name = $3, last_name = $4,
    birth_year = $5, birth_month = $6, birth_day = $7, birth_precision = $8,
    phone = $9, address_id = $10, updated_at = now()
WHERE national_id_hash = $1
RETURNING id`

	var id uuid.UUID
	err := tx.QueryRow(ctx, insert, p.NationalIDHash, p.NationalIDEnc, p.Title, p.FirstName, p.LastName,
		p.BirthYear, p.BirthMonth, p.BirthDay, p.BirthPrecision, p.Phone, p.AddressID).Scan(&id)
	if err == nil {
		return id, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, err
	}

	err = tx.QueryRow(ctx, update, p.NationalIDHash, p.Title, p.FirstName, p.LastName,
		p.BirthYear, p.BirthMonth, p.BirthDay, p.BirthPrecision, p.Phone, p.AddressID).Scan(&id)
	return id, err
}

func (r *Repo) FindCitizenIDByHash(ctx context.Context, hash string) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.pool.QueryRow(ctx, `SELECT id FROM citizens WHERE national_id_hash = $1`, hash).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, ErrCitizenNotFound
	}
	return id, err
}
