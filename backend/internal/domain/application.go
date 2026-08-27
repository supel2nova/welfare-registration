package domain

type Status string

const (
	StatusSubmitted Status = "SUBMITTED"
	StatusCancelled Status = "CANCELLED"
)

type SubmissionChannel string

const (
	ChannelWalkIn      SubmissionChannel = "WALK_IN"
	ChannelOnline      SubmissionChannel = "ONLINE"
	ChannelPartnerAPI  SubmissionChannel = "PARTNER_API"
	ChannelBatchImport SubmissionChannel = "BATCH_IMPORT"
)

type IDVerifyMethod string

const (
	VerifyLaserCode       IDVerifyMethod = "LASER_CODE"
	VerifyChipRead        IDVerifyMethod = "CHIP_READ"
	VerifyManualCardCheck IDVerifyMethod = "MANUAL_CARD_CHECK"
)

type VerificationMethod string

const (
	MethodLaserCode           VerificationMethod = "LASER_CODE"
	MethodChipRead            VerificationMethod = "CHIP_READ"
	MethodThaID               VerificationMethod = "THAID"
	MethodManual              VerificationMethod = "MANUAL"
	MethodDeclared            VerificationMethod = "DECLARED"
	MethodPendingVerification VerificationMethod = "PENDING_VERIFICATION"
)

type Title string

const (
	TitleMr   Title = "นาย"
	TitleMrs  Title = "นาง"
	TitleMiss Title = "นางสาว"
)

type MaritalStatus string

const (
	MaritalSingle   MaritalStatus = "SINGLE"
	MaritalMarried  MaritalStatus = "MARRIED"
	MaritalDivorced MaritalStatus = "DIVORCED"
	MaritalWidowed  MaritalStatus = "WIDOWED"
)

type Relation string

const (
	RelationSpouse Relation = "SPOUSE"
	RelationChild  Relation = "CHILD"
	RelationParent Relation = "PARENT"
	RelationOther  Relation = "OTHER"
)

type IncomeSourceType string

const (
	IncomeSalary IncomeSourceType = "SALARY"
	IncomeAgri   IncomeSourceType = "AGRI"
	IncomeTrade  IncomeSourceType = "TRADE"
	IncomeRent   IncomeSourceType = "RENT"
	IncomeOther  IncomeSourceType = "OTHER"
)

type AssetType string

const (
	AssetDeposit          AssetType = "DEPOSIT"
	AssetLottery          AssetType = "LOTTERY"
	AssetBond             AssetType = "BOND"
	AssetSecurities       AssetType = "SECURITIES"
	AssetLandAgri         AssetType = "LAND_AGRI"
	AssetLandResidential  AssetType = "LAND_RESIDENTIAL"
	AssetCondo            AssetType = "CONDO"
	AssetVehicleCar       AssetType = "VEHICLE_CAR"
	AssetVehicleMotorbike AssetType = "VEHICLE_MOTORCYCLE"
	AssetVehicleTricycle  AssetType = "VEHICLE_TRICYCLE"
	AssetVehicleFarm      AssetType = "VEHICLE_FARM"
)

type Unit string

const (
	UnitTHB   Unit = "THB"
	UnitRai   Unit = "RAI"
	UnitSqWa  Unit = "SQ_WA"
	UnitSqM   Unit = "SQ_M"
	UnitCount Unit = "COUNT"
)

var assetUnit = map[AssetType]Unit{
	AssetDeposit:          UnitTHB,
	AssetLottery:          UnitTHB,
	AssetBond:             UnitTHB,
	AssetSecurities:       UnitTHB,
	AssetLandAgri:         UnitRai,
	AssetLandResidential:  UnitSqWa,
	AssetCondo:            UnitSqM,
	AssetVehicleCar:       UnitCount,
	AssetVehicleMotorbike: UnitCount,
	AssetVehicleTricycle:  UnitCount,
	AssetVehicleFarm:      UnitCount,
}

func UnitFor(t AssetType) (Unit, bool) {
	u, ok := assetUnit[t]
	return u, ok
}

type LiabilityType string

const (
	LoanHome     LiabilityType = "LOAN_HOME"
	LoanVehicle  LiabilityType = "LOAN_VEHICLE"
	LoanPersonal LiabilityType = "LOAN_PERSONAL"
	LoanAgri     LiabilityType = "LOAN_AGRI"
	LoanOther    LiabilityType = "OTHER"
)

const (
	MaxHouseholdMembers = 20
	MaxSpouse           = 1
	MaxIncomeSources    = 10
	MaxAssets           = 30
	MaxLiabilities      = 20
	MinFiscalYear       = 2020
	MinManualNoteLen    = 10
)

var (
	validTitle = map[Title]bool{
		TitleMr: true, TitleMrs: true, TitleMiss: true,
	}
	validMarital = map[MaritalStatus]bool{
		MaritalSingle: true, MaritalMarried: true, MaritalDivorced: true, MaritalWidowed: true,
	}
	validRelation = map[Relation]bool{
		RelationSpouse: true, RelationChild: true, RelationParent: true, RelationOther: true,
	}
	validIncomeType = map[IncomeSourceType]bool{
		IncomeSalary: true, IncomeAgri: true, IncomeTrade: true, IncomeRent: true, IncomeOther: true,
	}
	validLiabilityType = map[LiabilityType]bool{
		LoanHome: true, LoanVehicle: true, LoanPersonal: true, LoanAgri: true, LoanOther: true,
	}
	validIDVerifyMethod = map[IDVerifyMethod]bool{
		VerifyLaserCode: true, VerifyChipRead: true, VerifyManualCardCheck: true,
	}
)

func IsValidTitle(v Title) bool                   { return validTitle[v] }
func IsValidMaritalStatus(v MaritalStatus) bool   { return validMarital[v] }
func IsValidRelation(v Relation) bool             { return validRelation[v] }
func IsValidIncomeType(v IncomeSourceType) bool   { return validIncomeType[v] }
func IsValidLiabilityType(v LiabilityType) bool   { return validLiabilityType[v] }
func IsValidIDVerifyMethod(v IDVerifyMethod) bool { return validIDVerifyMethod[v] }
