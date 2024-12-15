package entities

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"regexp"
	"time"
)

type Admin struct {
	ID          uuid.UUID
	FamilyName  string
	Name        string
	Email       string
	Password    string
	PhoneNumber string
	IdentityId  string
	Birthday    time.Time
	Status      bool
	Role        bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type AdminUpdate struct {
	FamilyName  *string
	Name        *string
	Email       *string
	Password    *string
	PhoneNumber *string
	IdentityId  *string
	Birthday    *time.Time
	Status      *bool
	Role        *bool
}

func (ad *Admin) ValidateAdmin() error {
	return validation.ValidateStruct(ad,
		validation.Field(&ad.FamilyName, validation.Required, validation.Length(2, 255)),
		validation.Field(&ad.Name, validation.Required, validation.Length(2, 255)),
		validation.Field(&ad.Email, validation.Required, is.Email),
		validation.Field(&ad.Password, validation.Required, validation.Length(2, 255)),
		validation.Field(&ad.PhoneNumber, validation.Required, validation.Length(10, 14), validation.Match((regexp.MustCompile((`^\d+$`))))),
		validation.Field(&ad.IdentityId, validation.Required, validation.Length(10, 15), validation.Match((regexp.MustCompile((`^\d+$`))))),
		validation.Field(&ad.Birthday, validation.Required),
		validation.Field(&ad.Status, validation.Required),
		validation.Field(&ad.Role, validation.Required),
		validation.Field(&ad.CreatedAt, validation.Required),
		validation.Field(&ad.UpdatedAt, validation.Required, validation.Min(ad.CreatedAt)),
	)
}

func (ad *Admin) ValidateAdminUpdate() error {
	return validation.ValidateStruct(ad,
		validation.Field(&ad.FamilyName, validation.Length(2, 255)),
		validation.Field(&ad.Name, validation.Length(2, 255)),
		validation.Field(&ad.Email, is.Email),
		validation.Field(&ad.Password, validation.Length(2, 255)),
		validation.Field(&ad.PhoneNumber, validation.Length(10, 14), validation.Match((regexp.MustCompile((`^\d+$`))))),
		validation.Field(&ad.IdentityId, validation.Length(10, 15), validation.Match((regexp.MustCompile((`^\d+$`))))),
	)
}

func NewAdmin(
	familyName string,
	name string,
	email string,
	password string,
	phoneNumber string,
	identityId string,
	birthday time.Time,
) (*Admin, error) {
	admin := &Admin{
		ID:          uuid.New(),
		FamilyName:  familyName,
		Name:        name,
		Email:       email,
		Password:    password,
		PhoneNumber: phoneNumber,
		IdentityId:  identityId,
		Birthday:    birthday,
		Status:      true,
		Role:        false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := admin.ValidateAdmin(); err != nil {
		return nil, err
	}

	return admin, nil
}

func NewSuperAdmin(
	familyName string,
	name string,
	email string,
	password string,
	phoneNumber string,
	identityId string,
	birthday time.Time,
) (*Admin, error) {
	admin := &Admin{
		ID:          uuid.New(),
		FamilyName:  familyName,
		Name:        name,
		Email:       email,
		Password:    password,
		PhoneNumber: phoneNumber,
		IdentityId:  identityId,
		Birthday:    birthday,
		Status:      true,
		Role:        true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := admin.ValidateAdmin(); err != nil {
		return nil, err
	}

	return admin, nil
}
