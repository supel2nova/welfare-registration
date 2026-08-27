package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/supel2nova/welfare-registration/backend/internal/config"
	"github.com/supel2nova/welfare-registration/backend/internal/dto"
	"github.com/supel2nova/welfare-registration/backend/internal/handler"
	"github.com/supel2nova/welfare-registration/backend/internal/repository"
	"github.com/supel2nova/welfare-registration/backend/internal/service"
	"github.com/supel2nova/welfare-registration/backend/internal/testutil"
	"github.com/supel2nova/welfare-registration/backend/internal/verifier"
	"github.com/supel2nova/welfare-registration/backend/pkg/idcrypto"
)

const (
	testPepper = "dev-pepper-do-not-use-in-production"
	testEncKey = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	userBAAC   = "aaaaaaaa-0000-0000-0000-000000000002"
)

func s(v string) *string { return &v }
func i(v int) *int       { return &v }

func validBody() dto.CreateApplicationRequest {
	return dto.CreateApplicationRequest{
		FiscalYear: 2026,
		Personal: dto.Personal{
			NationalID:     "1234567890121",
			LaserID:        s("JT8-1234567-89"),
			IDVerifyMethod: "LASER_CODE",
			Title:          "นาย",
			FirstName:      "สมชาย",
			LastName:       "ใจดี",
			BirthYear:      1985,
			BirthMonth:     i(3),
			BirthDay:       i(12),
			BirthPrecision: "FULL",
			Phone:          "0812345678",
			IsFarmer:       true,
			Address: dto.Address{
				HouseNo:         "12/34",
				Moo:             s("5"),
				ProvinceCode:    "50",
				DistrictCode:    "5001",
				SubdistrictCode: "500101",
				PostalCode:      "50200",
			},
		},
		Financial: dto.Financial{
			IncomeSources:   []dto.IncomeSource{{SourceType: "AGRI", AnnualAmount: 45000}},
			ExpenseToOthers: 0,
			Assets: []dto.Asset{
				{AssetType: "DEPOSIT", Amount: json.Number("15000"), Unit: "THB", JointAccountHolders: i(1)},
			},
			Liabilities:   []dto.Liability{},
			HasCreditCard: false,
		},
	}
}

func testRouter(t *testing.T, cfg config.Config) http.Handler {
	t.Helper()
	pool := testutil.Pool(t)
	testutil.Reset(t, pool)

	cipher, err := idcrypto.New(testEncKey)
	if err != nil {
		t.Fatal(err)
	}
	repo := repository.New(pool)
	apps := service.NewApplicationService(repo, verifier.Stub{}, cipher, testPepper)

	return handler.NewRouter(handler.Deps{
		Cfg:  cfg,
		Pool: pool,
		Repo: repo,
		Apps: apps,
		Ref:  service.NewRefService(repo),
	})
}

func postJSON(t *testing.T, h http.Handler, path string, headers map[string]string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatal(err)
		}
	}
	req := httptest.NewRequest(http.MethodPost, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w
}

func TestStubAuth_DisabledInProduction(t *testing.T) {
	r := testRouter(t, config.Config{
		Env:               "production",
		CORSOrigin:        "http://localhost:5173",
		StubAuthEnabled:   false,
		StubDefaultUserID: userBAAC,
	})

	w := postJSON(t, r, "/api/v1/applications", map[string]string{
		"X-Debug-User-Id": userBAAC,
	}, validBody())

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d body = %s want 401", w.Code, w.Body.String())
	}
}

func TestStubAuth_UnknownUser(t *testing.T) {
	r := testRouter(t, config.Config{
		Env:             "development",
		CORSOrigin:      "http://localhost:5173",
		StubAuthEnabled: true,
	})

	w := postJSON(t, r, "/api/v1/applications", map[string]string{
		"X-Debug-User-Id": uuid.New().String(),
	}, validBody())

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d body = %s want 401", w.Code, w.Body.String())
	}
}

func TestCreateApplication_HTTP400(t *testing.T) {
	r := testRouter(t, config.Config{
		Env:               "development",
		CORSOrigin:        "http://localhost:5173",
		StubAuthEnabled:   true,
		StubDefaultUserID: userBAAC,
	})
	req := validBody()
	req.Personal.NationalID = "123"
	w := postJSON(t, r, "/api/v1/applications", map[string]string{"X-Debug-User-Id": userBAAC}, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body = %s", w.Code, w.Body.String())
	}
}

func TestCreateApplication_HTTP201(t *testing.T) {
	r := testRouter(t, config.Config{
		Env:               "development",
		CORSOrigin:        "http://localhost:5173",
		StubAuthEnabled:   true,
		StubDefaultUserID: userBAAC,
	})
	w := postJSON(t, r, "/api/v1/applications", map[string]string{"X-Debug-User-Id": userBAAC}, validBody())
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body = %s", w.Code, w.Body.String())
	}
	var env struct {
		Data dto.CreateApplicationResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatal(err)
	}
	if env.Data.ApplicationNo == "" || env.Data.Status != "SUBMITTED" {
		t.Fatalf("data = %+v", env.Data)
	}
}

func TestCreateApplication_HTTP409(t *testing.T) {
	r := testRouter(t, config.Config{
		Env:               "development",
		CORSOrigin:        "http://localhost:5173",
		StubAuthEnabled:   true,
		StubDefaultUserID: userBAAC,
	})
	headers := map[string]string{"X-Debug-User-Id": userBAAC}
	if w := postJSON(t, r, "/api/v1/applications", headers, validBody()); w.Code != http.StatusCreated {
		t.Fatalf("setup status = %d body = %s", w.Code, w.Body.String())
	}
	w := postJSON(t, r, "/api/v1/applications", headers, validBody())
	if w.Code != http.StatusConflict {
		t.Fatalf("status = %d body = %s", w.Code, w.Body.String())
	}
	var env map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatal(err)
	}
	if env["errorCode"] != "DUP001" || env["data"] == nil {
		t.Fatalf("body = %s", w.Body.String())
	}
}
