package entities

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"time"
)

type User struct {
	ID           uuid.UUID           `validate:"omitempty,uuid4"`
	FamilyName   string              `validate:"required,min=2"`
	Name         string              `validate:"required,min=2"`
	Email        string              `validate:"required,email"`
	Password     string              `validate:"required,min=8"`
	PhoneNumber  string              `validate:"required,min=10,max=14,numeric"`
	Birthday     time.Time           `validate:"required"`
	AvatarUrl    string              `validate:"omitempty,url"`
	CapwallUrl   string              `validate:"omitempty,url"`
	Privacy      consts.PrivacyLevel `validate:"omitempty,oneof=public private friend_only"`
	Biography    string              `validate:"omitempty,max=500"`
	AuthType     consts.AuthType     `validate:"omitempty,oneof=local google"`
	AuthGoogleId string              `validate:"omitempty"`
	PostCount    int                 `validate:"gte=0"`
	FriendCount  int                 `validate:"gte=0"`
	Status       bool                `validate:"required"`
	CreatedAt    time.Time           `validate:"required"`
	UpdatedAt    time.Time           `validate:"required,gtefield=CreatedAt"`
	Setting      *Setting            `validate:"omitempty"`
}

type UserUpdate struct {
	FamilyName   *string              `validate:"omitempty,min=2"`
	Name         *string              `validate:"omitempty,min=2"`
	Email        *string              `validate:"omitempty,email"`
	Password     *string              `validate:"omitempty,min=8"`
	PhoneNumber  *string              `validate:"omitempty,min=10,max=14,numeric"`
	Birthday     *time.Time           `validate:"omitempty,date"`
	AvatarUrl    *string              `validate:"omitempty,url"`
	CapwallUrl   *string              `validate:"omitempty,url"`
	Privacy      *consts.PrivacyLevel `validate:"omitempty,oneof=public private friend_only"`
	Biography    *string              `validate:"omitempty,max=500"`
	AuthType     *consts.AuthType     `validate:"omitempty,oneof=local google"`
	AuthGoogleId *string              `validate:"omitempty"`
	PostCount    *int                 `validate:"omitempty,gte=0"`
	FriendCount  *int                 `validate:"omitempty,gte=0"`
	Status       *bool                `validate:"omitempty"`
	UpdatedAt    *time.Time           `validate:"omitempty,gtefield=CreatedAt"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *UserUpdate) ValidateUserUpdate() error {
	validate := validator.New()
	return validate.Struct(u)
}

func NewUser(
	familyName string,
	name string,
	email string,
	password string,
	phoneNumber string,
	birthday time.Time,
	authType consts.AuthType,
) (*User, error) {
	user := &User{
		ID:          uuid.New(),
		FamilyName:  familyName,
		Name:        name,
		Email:       email,
		Password:    password,
		PhoneNumber: phoneNumber,
		Birthday:    birthday,
		AvatarUrl:   consts.AVATAR_URL,
		CapwallUrl:  consts.CAPWALL_URL,
		Privacy:     consts.PUBLIC,
		AuthType:    authType,
		PostCount:   0,
		FriendCount: 0,
		Status:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

func NewUserUpdate(
	updateData *UserUpdate,
) (*UserUpdate, error) {
	userUpdate := &UserUpdate{
		FamilyName:   updateData.FamilyName,
		Name:         updateData.Name,
		Email:        updateData.Email,
		Password:     updateData.Password,
		PhoneNumber:  updateData.PhoneNumber,
		Birthday:     updateData.Birthday,
		AvatarUrl:    updateData.AvatarUrl,
		CapwallUrl:   updateData.CapwallUrl,
		Privacy:      updateData.Privacy,
		Biography:    updateData.Biography,
		AuthType:     updateData.AuthType,
		AuthGoogleId: updateData.AuthGoogleId,
		PostCount:    updateData.PostCount,
		FriendCount:  updateData.FriendCount,
		Status:       updateData.Status,
		UpdatedAt:    updateData.UpdatedAt,
	}

	if err := userUpdate.ValidateUserUpdate(); err != nil {
		return nil, err
	}

	return userUpdate, nil
}
