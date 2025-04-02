package entities

import (
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
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
	Password    *string
	PhoneNumber *string
	IdentityId  *string
	Birthday    *time.Time
	Status      *bool
	Role        *bool
}

func (ad *Admin) ValidateAdmin() error {
	return validation.ValidateStruct(ad,
		validation.Field(&ad.FamilyName, validation.Required, validation.RuneLength(2, 255)),
		validation.Field(&ad.Name, validation.Required, validation.RuneLength(2, 255)),
		validation.Field(&ad.Email, validation.Required, is.Email),
		validation.Field(&ad.Password, validation.Required, validation.RuneLength(2, 255)),
		validation.Field(&ad.PhoneNumber, validation.Required, validation.RuneLength(10, 14), validation.Match((regexp.MustCompile((`^\d+$`))))),
		validation.Field(&ad.IdentityId, validation.Required, validation.RuneLength(10, 15), validation.Match((regexp.MustCompile((`^\d+$`))))),
		validation.Field(&ad.Birthday, validation.Required),
		validation.Field(&ad.Status, validation.Required),
		validation.Field(&ad.CreatedAt, validation.Required),
		validation.Field(&ad.UpdatedAt, validation.Required, validation.Min(ad.CreatedAt)),
	)
}

func (ad *AdminUpdate) ValidateAdminUpdate() error {
	return validation.ValidateStruct(ad,
		validation.Field(&ad.FamilyName, validation.RuneLength(2, 255)),
		validation.Field(&ad.Name, validation.RuneLength(2, 255)),
		validation.Field(&ad.Password, validation.RuneLength(2, 255)),
		validation.Field(&ad.PhoneNumber, validation.RuneLength(10, 14), validation.Match((regexp.MustCompile((`^\d+$`))))),
		validation.Field(&ad.IdentityId, validation.RuneLength(10, 15), validation.Match((regexp.MustCompile((`^\d+$`))))),
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
	role bool,
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
		Role:        role,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := admin.ValidateAdmin(); err != nil {
		return nil, err
	}

	return admin, nil
}
