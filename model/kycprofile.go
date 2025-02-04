package models

import (
	"encoding/json"

	"salon_be/common"
)

const KYCProfileEntityName = "KYC Profile"

type IDDocument struct {
	common.SQLModel    `json:",inline"`
	KYCProfileID       uint32          `json:"-" gorm:"column:kyc_profile_id"`
	Type               int             `json:"type" gorm:"column:type"`
	Name               string          `json:"name" gorm:"column:name"`
	CardType           string          `json:"card_type" gorm:"column:card_type"`
	ID                 string          `json:"id" gorm:"column:id"`
	IDProbs            string          `json:"id_probs" gorm:"column:id_probs"`
	BirthDay           string          `json:"birth_day" gorm:"column:birth_day"`
	BirthDayProb       float64         `json:"birth_day_prob" gorm:"column:birth_day_prob"`
	Nationality        string          `json:"nationality" gorm:"column:nationality"`
	Gender             string          `json:"gender" gorm:"column:gender"`
	GenderProb         float64         `json:"gender_prob" gorm:"column:gender_prob"`
	ValidDate          string          `json:"valid_date" gorm:"column:valid_date"`
	ValidDateProb      float64         `json:"valid_date_prob" gorm:"column:valid_date_prob"`
	IssueDate          string          `json:"issue_date" gorm:"column:issue_date"`
	IssueDateProb      float64         `json:"issue_date_prob" gorm:"column:issue_date_prob"`
	IssuePlace         string          `json:"issue_place" gorm:"column:issue_place"`
	IssuePlaceProb     float64         `json:"issue_place_prob" gorm:"column:issue_place_prob"`
	OriginLocation     string          `json:"origin_location" gorm:"column:origin_location"`
	OriginLocationProb float64         `json:"origin_location_prob" gorm:"column:origin_location_prob"`
	RecentLocation     string          `json:"recent_location" gorm:"column:recent_location"`
	RecentLocationProb float64         `json:"recent_location_prob" gorm:"column:recent_location_prob"`
	PostCode           json.RawMessage `json:"post_code" gorm:"type:json"`
	Tampering          json.RawMessage `json:"tampering" gorm:"type:json"`
	ExpireWarning      string          `json:"expire_warning" gorm:"column:expire_warning"`
	IDFakeWarning      string          `json:"id_fake_warning" gorm:"column:id_fake_warning"`
	IDFakeProb         float64         `json:"id_fake_prob" gorm:"column:id_fake_prob"`
	LivenessStatus     string          `json:"liveness" gorm:"column:liveness"`
	LivenessMsg        string          `json:"liveness_msg" gorm:"column:liveness_msg"`
	FaceSwapping       bool            `json:"face_swapping" gorm:"column:face_swapping"`
	FakeLiveness       bool            `json:"fake_liveness" gorm:"column:fake_liveness"`
}

func (IDDocument) TableName() string {
	return "id_documents"
}

func (d *IDDocument) IsDocumentValid() bool {
	return d.LivenessStatus == "success" && !d.FakeLiveness && !d.FaceSwapping
}

type FaceVerification struct {
	common.SQLModel `json:",inline"`
	KYCProfileID    uint32  `json:"-" gorm:"column:kyc_profile_id"`
	Result          string  `json:"result" gorm:"column:result"`
	Msg             string  `json:"msg" gorm:"column:msg"`
	Prob            float64 `json:"prob" gorm:"column:prob"`
	LivenessStatus  string  `json:"liveness" gorm:"column:liveness"`
	LivenessMsg     string  `json:"liveness_msg" gorm:"column:liveness_msg"`
	IsEyeOpen       string  `json:"is_eye_open" gorm:"column:is_eye_open"`
	Masked          string  `json:"masked" gorm:"column:masked"`
}

func (FaceVerification) TableName() string {
	return "face_verifications"
}

func (f *FaceVerification) IsLivenessValid() bool {
	return f.LivenessStatus == "success"
}

func (f *FaceVerification) IsFaceMatch() bool {
	return f.Msg == "MATCH"
}

func (f *FaceVerification) HasMask() bool {
	return f.Masked == "yes"
}

type KYCProfile struct {
	common.SQLModel `json:",inline"`
	UserID          uint32            `json:"user_id" gorm:"column:user_id;uniqueIndex"`
	User            *User             `json:"-" gorm:"foreignKey:UserID"`
	CardID          string            `json:"card_id" gorm:"column:card_id"`
	PassportID      string            `json:"passport_id" gorm:"column:passport_id"`
	DriverLicenseID string            `json:"driver_license_id" gorm:"column:driver_license_id"`
	MilitaryID      string            `json:"military_id" gorm:"column:military_id"`
	PoliceID        string            `json:"police_id" gorm:"column:police_id"`
	OtherID         string            `json:"other_id" gorm:"column:other_id"`
	FullName        string            `json:"fullname" gorm:"column:fullname"`
	DOB             string            `json:"dob" gorm:"column:dob"`
	Gender          string            `json:"gender" gorm:"column:gender"`
	Address         string            `json:"address" gorm:"column:address"`
	Hometown        string            `json:"hometown" gorm:"column:hometown"`
	Nationality     string            `json:"nationality" gorm:"column:nationality"`
	IPFS            string            `json:"ipfs" gorm:"column:ipfs"`
	Title           string            `json:"title" gorm:"column:title"`
	OtherType       string            `json:"other_type" gorm:"column:other_type"`
	ExtraInfo       json.RawMessage   `json:"extra_info" gorm:"type:json"`
	Documents       []IDDocument      `json:"-" gorm:"foreignKey:KYCProfileID"`
	FaceData        *FaceVerification `json:"-" gorm:"foreignKey:KYCProfileID"`
	ClientSession   string            `json:"-" gorm:"column:client_session"`
}

func (KYCProfile) TableName() string {
	return "kyc_profiles"
}

func (k *KYCProfile) HasValidDocument() bool {
	if len(k.Documents) == 0 {
		return false
	}

	for _, doc := range k.Documents {
		if doc.LivenessStatus == "success" && !doc.FakeLiveness {
			return true
		}
	}
	return false
}
