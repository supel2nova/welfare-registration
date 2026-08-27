package service

import (
	"context"
	"errors"
	"regexp"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/supel2nova/welfare-registration/backend/internal/domain"
	"github.com/supel2nova/welfare-registration/backend/internal/dto"
	"github.com/supel2nova/welfare-registration/backend/internal/repository"
	"github.com/supel2nova/welfare-registration/backend/internal/testutil"
	"github.com/supel2nova/welfare-registration/backend/internal/verifier"
	"github.com/supel2nova/welfare-registration/backend/pkg/apperror"
	"github.com/supel2nova/welfare-registration/backend/pkg/idcrypto"
)

const (
	testPepper = "dev-pepper-do-not-use-in-production"
	testEncKey = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
)

var (
	orgBAAC = uuid.MustParse("22222222-2222-2222-2222-222222222222")
	orgKTB  = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	userBAAC = uuid.MustParse("aaaaaaaa-0000-0000-0000-000000000002")
	userKTB  = uuid.MustParse("aaaaaaaa-0000-0000-0000-000000000001")
	appNoRE  = regexp.MustCompile(`^WC-2026-\d{7}$`)
)

func baacActor() domain.Actor {
	uid := userBAAC
	return domain.Actor{Type: domain.ActorUser, OrgID: orgBAAC, UserID: &uid, Role: domain.RoleRegistrar}
}

func ktbActor() domain.Actor {
	uid := userKTB
	return domain.Actor{Type: domain.ActorUser, OrgID: orgKTB, UserID: &uid, Role: domain.RoleRegistrar}
}

func newAppService(t *testing.T, v verifier.Verifier) *ApplicationService {
	t.Helper()
	pool := testutil.Pool(t)
	testutil.Reset(t, pool)
	cipher, err := idcrypto.New(testEncKey)
	if err != nil {
		t.Fatal(err)
	}
	if v == nil {
		v = verifier.Stub{}
	}
	svc := NewApplicationService(repository.New(pool), v, cipher, testPepper)
	svc.now = func() time.Time { return today }
	return svc
}

func asAppErr(t *testing.T, err error) *apperror.Error {
	t.Helper()
	var ae *apperror.Error
	if !errors.As(err, &ae) {
		t.Fatalf("err = %v want *apperror.Error", err)
	}
	return ae
}

func count(t *testing.T, svc *ApplicationService, table string) int {
	t.Helper()
	var n int
	if err := svc.repo.Pool().QueryRow(context.Background(), "SELECT count(*) FROM "+table).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

func TestCreate_Success(t *testing.T) {
	svc := newAppService(t, nil)
	ctx := context.Background()

	res, err := svc.Create(ctx, baacActor(), validRequest())
	if err != nil {
		t.Fatal(err)
	}
	if !appNoRE.MatchString(res.ApplicationNo) {
		t.Fatalf("application_no = %q", res.ApplicationNo)
	}
	if res.Status != string(domain.StatusSubmitted) {
		t.Fatalf("status = %q", res.Status)
	}
	if res.RegistrationUnit != "ธ.ก.ส. สาขาแม่ริม" {
		t.Fatalf("unit = %q", res.RegistrationUnit)
	}
	if count(t, svc, "citizens") != 1 || count(t, svc, "applications") != 1 {
		t.Fatalf("citizens=%d applications=%d", count(t, svc, "citizens"), count(t, svc, "applications"))
	}
}

func TestCreate_AddressMismatch(t *testing.T) {
	svc := newAppService(t, nil)
	req := validRequest()
	req.Personal.Address.SubdistrictCode = "500701"

	_, err := svc.Create(context.Background(), baacActor(), req)
	ae := asAppErr(t, err)
	if ae.HTTPStatus != 400 || ae.Code != apperror.CodeAddress {
		t.Fatalf("got %+v", ae)
	}
}

func TestCreate_OrgFromActorNotBody(t *testing.T) {
	svc := newAppService(t, nil)
	ctx := context.Background()

	res, err := svc.Create(ctx, ktbActor(), validRequest())
	if err != nil {
		t.Fatal(err)
	}
	if res.RegistrationUnit != "ธนาคารกรุงไทย สาขาสีลม" {
		t.Fatalf("unit = %q", res.RegistrationUnit)
	}

	var unitID uuid.UUID
	if err := svc.repo.Pool().QueryRow(ctx, `SELECT registration_unit_id FROM applications WHERE application_no = $1`, res.ApplicationNo).Scan(&unitID); err != nil {
		t.Fatal(err)
	}
	if unitID != orgKTB {
		t.Fatalf("registration_unit_id = %v want %v", unitID, orgKTB)
	}
}

func TestCreate_DuplicateReturns409(t *testing.T) {
	svc := newAppService(t, nil)
	ctx := context.Background()
	req := validRequest()

	if _, err := svc.Create(ctx, baacActor(), req); err != nil {
		t.Fatal(err)
	}

	_, err := svc.Create(ctx, baacActor(), req)
	ae := asAppErr(t, err)
	if ae.HTTPStatus != 409 || ae.Code != apperror.CodeDuplicate {
		t.Fatalf("got HTTP %d code %s want 409 DUP001", ae.HTTPStatus, ae.Code)
	}
	dup, ok := ae.Data.(*dto.DuplicateInfo)
	if !ok || dup == nil {
		t.Fatalf("data = %#v", ae.Data)
	}
	if dup.RegisteredUnit != "ธ.ก.ส. สาขาแม่ริม" || dup.ApplicationNo == "" {
		t.Fatalf("dup = %+v", dup)
	}
}

func TestCreate_DuplicateSkipsProvider(t *testing.T) {
	spy := &verifier.Spy{Inner: verifier.Stub{}}
	svc := newAppService(t, spy)
	ctx := context.Background()
	req := validRequest()

	if _, err := svc.Create(ctx, baacActor(), req); err != nil {
		t.Fatal(err)
	}
	firstCalls := spy.Calls

	_, err := svc.Create(ctx, baacActor(), req)
	asAppErr(t, err)
	if spy.Calls != firstCalls {
		t.Fatalf("spy.Calls = %d เพิ่มจาก %d — ซ้ำต้องไม่เรียก provider", spy.Calls, firstCalls)
	}
}

func TestCreate_CancelledAllowsResubmit(t *testing.T) {
	svc := newAppService(t, nil)
	ctx := context.Background()
	req := validRequest()

	if _, err := svc.Create(ctx, baacActor(), req); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.repo.Pool().Exec(ctx, `UPDATE applications SET status = 'CANCELLED'`); err != nil {
		t.Fatal(err)
	}

	res, err := svc.Create(ctx, baacActor(), req)
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != string(domain.StatusSubmitted) {
		t.Fatalf("status = %q", res.Status)
	}
	if count(t, svc, "applications") != 2 {
		t.Fatalf("applications = %d want 2", count(t, svc, "applications"))
	}
}

func TestCreate_DifferentFiscalYear(t *testing.T) {
	svc := newAppService(t, nil)
	ctx := context.Background()
	req := validRequest()
	req.FiscalYear = 2025

	if _, err := svc.Create(ctx, baacActor(), req); err != nil {
		t.Fatal(err)
	}

	req.FiscalYear = 2026
	if _, err := svc.Create(ctx, baacActor(), req); err != nil {
		t.Fatal(err)
	}
	if count(t, svc, "applications") != 2 {
		t.Fatalf("applications = %d want 2", count(t, svc, "applications"))
	}
}

func TestCreate_ConcurrentSameNationalID(t *testing.T) {
	svc := newAppService(t, nil)
	ctx := context.Background()
	req := validRequest()

	const n = 10
	var wg sync.WaitGroup
	errs := make([]error, n)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			_, errs[i] = svc.Create(ctx, baacActor(), req)
		}(i)
	}
	wg.Wait()

	ok, conflict := 0, 0
	for _, err := range errs {
		if err == nil {
			ok++
			continue
		}
		ae := asAppErr(t, err)
		if ae.HTTPStatus == 409 {
			conflict++
			continue
		}
		t.Fatalf("unexpected err: %+v", ae)
	}
	if ok != 1 || conflict != 9 {
		t.Fatalf("ok=%d conflict=%d want 1/9", ok, conflict)
	}
	if count(t, svc, "citizens") != 1 || count(t, svc, "applications") != 1 {
		t.Fatalf("citizens=%d applications=%d", count(t, svc, "citizens"), count(t, svc, "applications"))
	}
}

func TestCreate_LaserIDNeverStored(t *testing.T) {
	svc := newAppService(t, nil)
	ctx := context.Background()
	req := validRequest()
	laser := *req.Personal.LaserID

	if _, err := svc.Create(ctx, baacActor(), req); err != nil {
		t.Fatal(err)
	}

	tables := []string{
		"citizens", "addresses", "applications", "identity_verifications",
		"household_members", "income_sources", "assets", "liabilities",
		"application_status_history",
	}
	laserRE := regexp.MustCompile(regexp.QuoteMeta(laser))
	keyRE := regexp.MustCompile(`"laser_id"`)
	for _, table := range tables {
		rows, err := svc.repo.Pool().Query(ctx, `SELECT CAST(row_to_json(t) AS text) FROM `+table+` t`)
		if err != nil {
			t.Fatal(err)
		}
		for rows.Next() {
			var raw string
			if err := rows.Scan(&raw); err != nil {
				rows.Close()
				t.Fatal(err)
			}
			if laserRE.MatchString(raw) {
				rows.Close()
				t.Fatalf("พบ laser_id ในตาราง %s: %s", table, raw)
			}
			if keyRE.MatchString(raw) {
				rows.Close()
				t.Fatalf("พบ key laser_id ในตาราง %s: %s", table, raw)
			}
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			t.Fatal(err)
		}
	}

	var snap string
	if err := svc.repo.Pool().QueryRow(ctx, `SELECT applicant_snapshot::text FROM applications`).Scan(&snap); err != nil {
		t.Fatal(err)
	}
	if regexp.MustCompile(`laser_id`).MatchString(snap) {
		t.Fatalf("snapshot มี laser_id: %s", snap)
	}
}

func TestCreate_ProviderUnavailable(t *testing.T) {
	svc := newAppService(t, verifier.Stub{Unavailable: true})
	ctx := context.Background()

	res, err := svc.Create(ctx, baacActor(), validRequest())
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != string(domain.StatusSubmitted) {
		t.Fatalf("status = %q", res.Status)
	}

	var method string
	if err := svc.repo.Pool().QueryRow(ctx, `SELECT method FROM identity_verifications`).Scan(&method); err != nil {
		t.Fatal(err)
	}
	if method != string(domain.MethodPendingVerification) {
		t.Fatalf("method = %q", method)
	}
}

func TestCreate_KYCMismatch(t *testing.T) {
	svc := newAppService(t, verifier.Stub{FailAll: true})
	ctx := context.Background()

	_, err := svc.Create(ctx, baacActor(), validRequest())
	ae := asAppErr(t, err)
	if ae.HTTPStatus != 400 || ae.Code != apperror.CodeKYC {
		t.Fatalf("got %+v", ae)
	}
	if count(t, svc, "citizens") != 0 || count(t, svc, "applications") != 0 {
		t.Fatal("KYC ไม่ผ่านต้องไม่เขียน DB")
	}
}

func TestCreate_ManualSkipsProvider(t *testing.T) {
	spy := &verifier.Spy{Inner: verifier.Stub{}}
	svc := newAppService(t, spy)
	ctx := context.Background()

	req := validRequest()
	req.Personal.LaserID = nil
	req.Personal.IDVerifyMethod = "MANUAL_CARD_CHECK"
	req.Personal.IDVerifyNote = s("บัตรตลอดชีพ อ่านรหัสหลังบัตรไม่ได้")

	res, err := svc.Create(ctx, baacActor(), req)
	if err != nil {
		t.Fatal(err)
	}
	if spy.Calls != 0 {
		t.Fatalf("spy.Calls = %d want 0", spy.Calls)
	}

	var method string
	if err := svc.repo.Pool().QueryRow(ctx, `SELECT method FROM identity_verifications WHERE citizen_id = (
		SELECT citizen_id FROM applications WHERE application_no = $1)`, res.ApplicationNo).Scan(&method); err != nil {
		t.Fatal(err)
	}
	if method != string(domain.MethodManual) {
		t.Fatalf("method = %q", method)
	}
}

func TestCreate_RollbackOnChildrenFailure(t *testing.T) {
	svc := newAppService(t, nil)
	svc.insertChildren = func(ctx context.Context, tx pgx.Tx, appID uuid.UUID, members []repository.MemberParams, f dto.Financial) error {
		return errors.New("mock insert children failed")
	}
	ctx := context.Background()

	_, err := svc.Create(ctx, baacActor(), validRequest())
	ae := asAppErr(t, err)
	if ae.HTTPStatus != 500 {
		t.Fatalf("got HTTP %d want 500", ae.HTTPStatus)
	}
	if count(t, svc, "citizens") != 0 || count(t, svc, "applications") != 0 {
		t.Fatalf("citizens=%d applications=%d ต้อง rollback หมด", count(t, svc, "citizens"), count(t, svc, "applications"))
	}
}

func TestBuildSnapshot_NoLaserID(t *testing.T) {
	req := validRequest()
	raw, err := buildSnapshot(req, repository.ResolvedAddress{
		SubdistrictName: "ศรีภูมิ", DistrictName: "เมืองเชียงใหม่", ProvinceName: "เชียงใหม่",
	})
	if err != nil {
		t.Fatal(err)
	}
	if regexp.MustCompile(`laser_id`).MatchString(string(raw)) {
		t.Fatalf("snapshot มี laser_id: %s", raw)
	}
	if !regexp.MustCompile(`national_id_mask`).MatchString(string(raw)) {
		t.Fatal("ต้องมี national_id_mask")
	}
}
