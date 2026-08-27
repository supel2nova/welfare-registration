package domain_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/supel2nova/welfare-registration/backend/internal/domain"
)

func TestActor(t *testing.T) {
	userID := uuid.MustParse("aaaaaaaa-0000-0000-0000-000000000002")
	clientID := uuid.MustParse("bbbbbbbb-0000-0000-0000-000000000001")
	orgID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	cases := []struct {
		name     string
		actor    domain.Actor
		isSystem bool
		actorID  uuid.UUID
	}{
		{"เจ้าหน้าที่ล็อกอิน", domain.Actor{Type: domain.ActorUser, OrgID: orgID, UserID: &userID, Role: domain.RoleRegistrar}, false, userID},
		{"ระบบพาร์ตเนอร์ยิง API", domain.Actor{Type: domain.ActorSystem, OrgID: orgID, ClientID: &clientID, Role: domain.RoleRegistrar}, true, clientID},
		{"USER แต่ไม่มี user id (ไม่ควรเกิด) ไม่ panic", domain.Actor{Type: domain.ActorUser, OrgID: orgID}, false, uuid.Nil},
		{"SYSTEM แต่ไม่มี client id (ไม่ควรเกิด) ไม่ panic", domain.Actor{Type: domain.ActorSystem, OrgID: orgID}, true, uuid.Nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.actor.IsSystem(); got != c.isSystem {
				t.Errorf("IsSystem = %v want %v", got, c.isSystem)
			}
			if got := c.actor.ActorID(); got != c.actorID {
				t.Errorf("ActorID = %v want %v", got, c.actorID)
			}
		})
	}
}

func TestUnitFor(t *testing.T) {
	cases := []struct {
		name  string
		asset domain.AssetType
		want  domain.Unit
		ok    bool
	}{
		{"เงินฝากเป็นบาท", domain.AssetDeposit, domain.UnitTHB, true},
		{"สลากเป็นบาท", domain.AssetLottery, domain.UnitTHB, true},
		{"ที่ดินเกษตรเป็นไร่", domain.AssetLandAgri, domain.UnitRai, true},
		{"ที่ดินอยู่อาศัยเป็นตารางวา", domain.AssetLandResidential, domain.UnitSqWa, true},
		{"คอนโดเป็นตารางเมตร", domain.AssetCondo, domain.UnitSqM, true},
		{"รถยนต์นับเป็นคัน", domain.AssetVehicleCar, domain.UnitCount, true},
		{"รถไถนับเป็นคัน", domain.AssetVehicleFarm, domain.UnitCount, true},
		{"asset_type ที่ไม่รู้จัก", domain.AssetType("SPACESHIP"), "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, ok := domain.UnitFor(c.asset)
			if ok != c.ok || got != c.want {
				t.Errorf("UnitFor(%s) = %q,%v want %q,%v", c.asset, got, ok, c.want, c.ok)
			}
		})
	}
}

func TestEnumGuards(t *testing.T) {
	if !domain.IsValidTitle(domain.TitleMiss) || domain.IsValidTitle("Mr.") {
		t.Error("title")
	}
	if !domain.IsValidMaritalStatus(domain.MaritalMarried) || domain.IsValidMaritalStatus("SEPARATED") {
		t.Error("marital")
	}
	if !domain.IsValidRelation(domain.RelationSpouse) || domain.IsValidRelation("FRIEND") {
		t.Error("relation")
	}
	if !domain.IsValidIncomeType(domain.IncomeAgri) || domain.IsValidIncomeType("CRYPTO") {
		t.Error("income")
	}
	if !domain.IsValidLiabilityType(domain.LoanAgri) || domain.IsValidLiabilityType("LOAN_MOON") {
		t.Error("liability")
	}
	if !domain.IsValidIDVerifyMethod(domain.VerifyManualCardCheck) || domain.IsValidIDVerifyMethod("MANUAL") {
		t.Error("id_verify_method")
	}
}
