package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/supel2nova/welfare-registration/backend/internal/domain"
	"github.com/supel2nova/welfare-registration/backend/internal/dto"
	"github.com/supel2nova/welfare-registration/backend/internal/repository"
	"github.com/supel2nova/welfare-registration/backend/internal/testutil"
)

var (
	orgBAAC = uuid.MustParse("22222222-2222-2222-2222-222222222222")
	userID  = uuid.MustParse("aaaaaaaa-0000-0000-0000-000000000002")
)

func newRepo(t *testing.T) *repository.Repo {
	t.Helper()
	pool := testutil.Pool(t)
	testutil.Reset(t, pool)
	return repository.New(pool)
}

func TestResolveAddress(t *testing.T) {
	r := newRepo(t)
	ctx := context.Background()

	t.Run("รหัสสอดคล้องกันทั้งชุด", func(t *testing.T) {
		got, err := r.ResolveAddress(ctx, "50", "5001", "500101", "50200")
		if err != nil {
			t.Fatal(err)
		}
		if got.SubdistrictName != "ศรีภูมิ" || got.DistrictName != "เมืองเชียงใหม่" || got.ProvinceName != "เชียงใหม่" {
			t.Fatalf("ชื่อไม่ตรง: %+v", got)
		}
		if !got.PostalMatches {
			t.Error("รหัสไปรษณีย์ควรตรง")
		}
	})

	t.Run("ตำบลหนึ่งมีได้หลายรหัสไปรษณีย์", func(t *testing.T) {
		for _, postal := range []string{"50200", "50100"} {
			got, err := r.ResolveAddress(ctx, "50", "5001", "500108", postal)
			if err != nil || !got.PostalMatches {
				t.Errorf("%s: %+v %v", postal, got, err)
			}
		}
	})

	t.Run("รหัสไปรษณีย์ไม่ตรงตำบล แต่ยังคืนชื่อมาให้", func(t *testing.T) {
		got, err := r.ResolveAddress(ctx, "50", "5001", "500101", "10200")
		if err != nil {
			t.Fatal(err)
		}
		if got.PostalMatches {
			t.Error("รหัสไปรษณีย์กรุงเทพไม่ควรตรงกับตำบลในเชียงใหม่")
		}
		if got.SubdistrictName == "" {
			t.Error("ต้องรู้ว่าผิดที่รหัสไปรษณีย์ ไม่ใช่ผิดทั้งชุด")
		}
	})

	t.Run("ตำบลไม่ได้อยู่ในอำเภอที่อ้าง", func(t *testing.T) {
		_, err := r.ResolveAddress(ctx, "50", "5001", "500701", "50180")
		if !errors.Is(err, repository.ErrAddressNotFound) {
			t.Fatalf("err = %v want ErrAddressNotFound", err)
		}
	})

	t.Run("อำเภอไม่ได้อยู่ในจังหวัดที่อ้าง", func(t *testing.T) {
		_, err := r.ResolveAddress(ctx, "10", "5001", "500101", "50200")
		if !errors.Is(err, repository.ErrAddressNotFound) {
			t.Fatalf("err = %v want ErrAddressNotFound", err)
		}
	})

	t.Run("รหัสไม่มีในระบบ", func(t *testing.T) {
		_, err := r.ResolveAddress(ctx, "99", "9999", "999999", "99999")
		if !errors.Is(err, repository.ErrAddressNotFound) {
			t.Fatalf("err = %v want ErrAddressNotFound", err)
		}
	})
}

func TestRefLists(t *testing.T) {
	r := newRepo(t)
	ctx := context.Background()

	provinces, err := r.Provinces(ctx)
	if err != nil || len(provinces) != 3 {
		t.Fatalf("provinces = %d %v", len(provinces), err)
	}

	districts, err := r.Districts(ctx, "50")
	if err != nil || len(districts) == 0 {
		t.Fatalf("districts = %d %v", len(districts), err)
	}
	for _, d := range districts {
		if d.Code[:2] != "50" {
			t.Errorf("อำเภอ %s ไม่ได้อยู่จังหวัด 50", d.Code)
		}
	}

	subdistricts, err := r.Subdistricts(ctx, "5007")
	if err != nil || len(subdistricts) == 0 {
		t.Fatalf("subdistricts = %d %v", len(subdistricts), err)
	}

	empty, err := r.Districts(ctx, "99")
	if err != nil {
		t.Fatal(err)
	}
	if empty == nil {
		t.Error("ไม่มีข้อมูลต้องคืน [] ไม่ใช่ null")
	}
}

func TestFindUserByID(t *testing.T) {
	r := newRepo(t)
	ctx := context.Background()

	t.Run("ผู้ใช้ที่มีอยู่", func(t *testing.T) {
		u, err := r.FindUserByID(ctx, userID)
		if err != nil {
			t.Fatal(err)
		}
		if u.OrgID != orgBAAC || u.Role != string(domain.RoleRegistrar) {
			t.Fatalf("ได้ %+v", u)
		}
	})

	t.Run("uuid ที่ไม่มีในระบบ", func(t *testing.T) {
		_, err := r.FindUserByID(ctx, uuid.New())
		if !errors.Is(err, repository.ErrUserNotFound) {
			t.Fatalf("err = %v want ErrUserNotFound", err)
		}
	})
}

func TestUpsertCitizen(t *testing.T) {
	r := newRepo(t)
	ctx := context.Background()

	first := insertCitizen(t, r, "hash-upsert", "สมชาย")
	second := insertCitizen(t, r, "hash-upsert", "สมชายแก้ไขแล้ว")

	if first != second {
		t.Fatalf("hash เดิมต้องได้ citizen เดิม: %v vs %v", first, second)
	}

	var name string
	if err := r.Pool().QueryRow(ctx, `SELECT first_name FROM citizens WHERE id = $1`, first).Scan(&name); err != nil {
		t.Fatal(err)
	}
	if name != "สมชายแก้ไขแล้ว" {
		t.Errorf("ข้อมูลไม่ถูกอัปเดต: %q", name)
	}

	var count int
	if err := r.Pool().QueryRow(ctx, `SELECT count(*) FROM citizens`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("citizens = %d แถว want 1", count)
	}
}

func TestFindCitizenIDByHash(t *testing.T) {
	r := newRepo(t)
	id := insertCitizen(t, r, "hash-find", "สมหญิง")

	got, err := r.FindCitizenIDByHash(context.Background(), "hash-find")
	if err != nil || got != id {
		t.Fatalf("got %v %v want %v", got, err, id)
	}

	if _, err := r.FindCitizenIDByHash(context.Background(), "ไม่มีจริง"); !errors.Is(err, repository.ErrCitizenNotFound) {
		t.Fatalf("err = %v want ErrCitizenNotFound", err)
	}
}

func TestUniqueViolationOnSecondApplication(t *testing.T) {
	r := newRepo(t)
	ctx := context.Background()

	citizenID := insertCitizen(t, r, "hash-dup", "สมปอง")
	insertApplication(t, r, citizenID, 2026)

	tx, err := r.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)

	appNo, err := r.NextAppNo(ctx, tx, 2026)
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = r.InsertApplication(ctx, tx, appParams(appNo, citizenID, 2026))

	if !repository.IsUniqueViolation(err, repository.ConstraintActivePerYear) {
		t.Fatalf("err = %v ต้องเป็น unique violation ของ %s", err, repository.ConstraintActivePerYear)
	}
	if repository.IsUniqueViolation(err, "constraint_อื่น") {
		t.Error("ต้องแยกได้ว่าชนที่ constraint ไหน")
	}
	if repository.IsUniqueViolation(errors.New("boom"), "") {
		t.Error("error ทั่วไปต้องไม่ถูกนับเป็น unique violation")
	}
}

func TestFindActiveByHash(t *testing.T) {
	r := newRepo(t)
	ctx := context.Background()

	citizenID := insertCitizen(t, r, "hash-active", "สมศรี")
	insertApplication(t, r, citizenID, 2026)

	t.Run("ปีเดียวกันเจอ", func(t *testing.T) {
		got, err := r.FindActiveByHash(ctx, "hash-active", 2026)
		if err != nil {
			t.Fatal(err)
		}
		if got == nil {
			t.Fatal("ต้องเจอใบเดิม")
		}
		if got.RegisteredUnit != "ธ.ก.ส. สาขาแม่ริม" {
			t.Errorf("registered_unit = %q", got.RegisteredUnit)
		}
		if len(got.RegisteredAt) != 10 {
			t.Errorf("registered_at ต้องเป็น YYYY-MM-DD ได้ %q", got.RegisteredAt)
		}
		if got.CanAppeal {
			t.Error("MVP1 ยังไม่มีอุทธรณ์")
		}
	})

	t.Run("คนละปีงบประมาณไม่เจอ", func(t *testing.T) {
		got, err := r.FindActiveByHash(ctx, "hash-active", 2025)
		if err != nil || got != nil {
			t.Fatalf("got %+v %v", got, err)
		}
	})

	t.Run("ใบที่ยกเลิกแล้วไม่นับ", func(t *testing.T) {
		if _, err := r.Pool().Exec(ctx, `UPDATE applications SET status = 'CANCELLED'`); err != nil {
			t.Fatal(err)
		}
		got, err := r.FindActiveByHash(ctx, "hash-active", 2026)
		if err != nil || got != nil {
			t.Fatalf("got %+v %v", got, err)
		}
	})
}

func TestInsertChildren(t *testing.T) {
	r := newRepo(t)
	ctx := context.Background()

	citizenID := insertCitizen(t, r, "hash-children", "สมหมาย")
	appID := insertApplication(t, r, citizenID, 2026)

	tx, err := r.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)

	hash := "hash-spouse"
	err = r.InsertChildren(ctx, tx, appID,
		[]repository.MemberParams{{Relation: "SPOUSE", NationalIDHash: &hash, FullName: "นางสมหญิง ใจดี"}},
		dto.Financial{
			IncomeSources: []dto.IncomeSource{{SourceType: "AGRI", AnnualAmount: 45000}},
			Assets: []dto.Asset{
				{AssetType: "DEPOSIT", Amount: "15000", Unit: "THB"},
				{AssetType: "LAND_AGRI", Amount: "8.50", Unit: "RAI"},
			},
			Liabilities: []dto.Liability{{LiabilityType: "LOAN_PERSONAL", CreditLimit: 30000}},
		})
	if err != nil {
		t.Fatal(err)
	}

	var (
		area   string
		joint  int
		income int64
	)
	if err := tx.QueryRow(ctx, `SELECT amount::text FROM assets WHERE asset_type = 'LAND_AGRI'`).Scan(&area); err != nil {
		t.Fatal(err)
	}
	if area != "8.50" {
		t.Errorf("ที่ดิน 8.5 ไร่เก็บเป็น %q", area)
	}
	if err := tx.QueryRow(ctx, `SELECT joint_account_holders FROM assets WHERE asset_type = 'DEPOSIT'`).Scan(&joint); err != nil {
		t.Fatal(err)
	}
	if joint != 1 {
		t.Errorf("joint_account_holders ที่ไม่ได้ระบุต้องเป็น 1 ได้ %d", joint)
	}
	if err := tx.QueryRow(ctx, `SELECT annual_amount FROM income_sources`).Scan(&income); err != nil {
		t.Fatal(err)
	}
	if income != 45000 {
		t.Errorf("annual_amount = %d", income)
	}
}

func insertCitizen(t *testing.T, r *repository.Repo, hash, firstName string) uuid.UUID {
	t.Helper()
	ctx := context.Background()

	tx, err := r.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)

	addrID, err := r.InsertAddress(ctx, tx, repository.AddressParams{
		HouseNo: "12/34", SubdistrictCode: "500701", DistrictCode: "5007", ProvinceCode: "50",
		PostalCode: "50180", SubdistrictName: "ริมใต้", DistrictName: "แม่ริม", ProvinceName: "เชียงใหม่",
	})
	if err != nil {
		t.Fatal(err)
	}

	id, err := r.UpsertCitizen(ctx, tx, repository.CitizenParams{
		NationalIDHash: hash, NationalIDEnc: []byte{0x01, 0x02}, Title: "นาย",
		FirstName: firstName, LastName: "ใจดี", BirthYear: 1985, BirthPrecision: "YEAR_ONLY",
		Phone: "0812345678", AddressID: addrID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}
	return id
}

func appParams(appNo string, citizenID uuid.UUID, fiscalYear int) repository.ApplicationParams {
	return repository.ApplicationParams{
		ApplicationNo:      appNo,
		CitizenID:          citizenID,
		FiscalYear:         fiscalYear,
		Snapshot:           []byte(`{"national_id_mask":"1-2345-xxxxx-xx-1"}`),
		RegistrationUnitID: orgBAAC,
		CreatedByUserID:    &userID,
		SubmissionChannel:  string(domain.ChannelWalkIn),
	}
}

func insertApplication(t *testing.T, r *repository.Repo, citizenID uuid.UUID, fiscalYear int) uuid.UUID {
	t.Helper()
	ctx := context.Background()

	tx, err := r.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)

	appNo, err := r.NextAppNo(ctx, tx, fiscalYear)
	if err != nil {
		t.Fatal(err)
	}
	id, _, err := r.InsertApplication(ctx, tx, appParams(appNo, citizenID, fiscalYear))
	if err != nil {
		t.Fatal(err)
	}
	if err := r.InsertStatusHistory(ctx, tx, id, domain.StatusSubmitted, domain.Actor{
		Type: domain.ActorUser, OrgID: orgBAAC, UserID: &userID, Role: domain.RoleRegistrar,
	}); err != nil {
		t.Fatal(err)
	}
	if err := tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}
	return id
}
