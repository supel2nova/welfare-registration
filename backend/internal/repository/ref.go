package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/supel2nova/welfare-registration/backend/internal/dto"
)

var ErrAddressNotFound = errors.New("repository: ไม่พบที่อยู่ตามรหัสที่ส่งมา")

type ResolvedAddress struct {
	SubdistrictName string
	DistrictName    string
	ProvinceName    string
	PostalMatches   bool
}

func (r *Repo) ResolveAddress(ctx context.Context, provinceCode, districtCode, subdistrictCode, postalCode string) (ResolvedAddress, error) {
	const q = `
SELECT s.name_th, d.name_th, p.name_th,
       EXISTS (SELECT 1 FROM ref_subdistrict_postal sp
               WHERE sp.subdistrict_code = s.code AND sp.postal_code = $4)
FROM ref_subdistricts s
JOIN ref_districts  d ON d.code = s.district_code
JOIN ref_provinces  p ON p.code = d.province_code
WHERE s.code = $1 AND d.code = $2 AND p.code = $3`

	var a ResolvedAddress
	err := r.pool.QueryRow(ctx, q, subdistrictCode, districtCode, provinceCode, postalCode).
		Scan(&a.SubdistrictName, &a.DistrictName, &a.ProvinceName, &a.PostalMatches)
	if errors.Is(err, pgx.ErrNoRows) {
		return a, ErrAddressNotFound
	}
	return a, err
}

const AddressSearchLimit = 30

func (r *Repo) SearchAddress(ctx context.Context, q string) ([]dto.AddressOption, error) {
	const sql = `
SELECT s.code, s.name_th, s.kind,
       d.code, d.name_th, d.kind,
       p.code, p.name_th,
       sp.postal_code
FROM ref_subdistricts s
JOIN ref_districts d ON d.code = s.district_code
JOIN ref_provinces p ON p.code = d.province_code
JOIN ref_subdistrict_postal sp ON sp.subdistrict_code = s.code
WHERE s.name_th ILIKE $1
   OR d.name_th ILIKE $1
   OR p.name_th ILIKE $1
   OR sp.postal_code LIKE $2
ORDER BY (s.name_th ILIKE $2) DESC, s.name_th, sp.is_primary DESC, sp.postal_code
LIMIT $3`

	rows, err := r.pool.Query(ctx, sql, "%"+q+"%", q+"%", AddressSearchLimit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []dto.AddressOption{}
	for rows.Next() {
		var a dto.AddressOption
		if err := rows.Scan(
			&a.SubdistrictCode, &a.SubdistrictName, &a.SubdistrictKind,
			&a.DistrictCode, &a.DistrictName, &a.DistrictKind,
			&a.ProvinceCode, &a.ProvinceName, &a.PostalCode,
		); err != nil {
			return nil, err
		}
		items = append(items, a)
	}
	return items, rows.Err()
}

func (r *Repo) Provinces(ctx context.Context) ([]dto.RefItem, error) {
	return r.refItems(ctx, `SELECT code, name_th FROM ref_provinces ORDER BY code`)
}

func (r *Repo) Districts(ctx context.Context, provinceCode string) ([]dto.RefItem, error) {
	return r.refItems(ctx, `SELECT code, name_th FROM ref_districts WHERE province_code = $1 ORDER BY code`, provinceCode)
}

func (r *Repo) Subdistricts(ctx context.Context, districtCode string) ([]dto.RefItem, error) {
	return r.refItems(ctx, `SELECT code, name_th FROM ref_subdistricts WHERE district_code = $1 ORDER BY code`, districtCode)
}

func (r *Repo) refItems(ctx context.Context, q string, args ...any) ([]dto.RefItem, error) {
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []dto.RefItem{}
	for rows.Next() {
		var it dto.RefItem
		if err := rows.Scan(&it.Code, &it.NameTH); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, rows.Err()
}
