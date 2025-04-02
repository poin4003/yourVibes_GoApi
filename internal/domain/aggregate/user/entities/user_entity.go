package entities

import (
	"regexp"
	"time"

	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/pointer"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type User struct {
	ID           uuid.UUID
	FamilyName   string
	Name         string
	Email        string
	Password     *string
	PhoneNumber  *string
	Birthday     *time.Time
	AvatarUrl    string
	CapwallUrl   string
	Privacy      consts.PrivacyLevel
	Biography    string
	AuthType     consts.AuthType
	AuthGoogleId *string
	PostCount    int
	FriendCount  int
	Status       bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Setting      *Setting
}

type UserUpdate struct {
	FamilyName   *string
	Name         *string
	PhoneNumber  *string
	Birthday     *time.Time
	Password     *string
	AvatarUrl    *string
	CapwallUrl   *string
	Privacy      *consts.PrivacyLevel
	Biography    *string
	AuthType     *consts.AuthType
	AuthGoogleId *string
	PostCount    *int
	FriendCount  *int
	Status       *bool
	UpdatedAt    *time.Time
}

func (u *User) ValidateUser() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.FamilyName, validation.Required, validation.RuneLength(2, 255)),
		validation.Field(&u.Name, validation.Required, validation.RuneLength(2, 255)),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.RuneLength(8, 255)),
		validation.Field(&u.PhoneNumber, validation.Required, validation.RuneLength(10, 14), validation.Match((regexp.MustCompile((`^\d+$`))))),
		validation.Field(&u.Birthday, validation.Required),
		validation.Field(&u.AvatarUrl, is.URL),
		validation.Field(&u.CapwallUrl, is.URL),
		validation.Field(&u.Privacy, validation.In(consts.PrivacyLevels...)),
		validation.Field(&u.Biography, validation.RuneLength(0, 500)),
		validation.Field(&u.AuthType, validation.In(consts.AuthTypes...)),
		validation.Field(&u.PostCount, validation.Min(0)),
		validation.Field(&u.FriendCount, validation.Min(0)),
		validation.Field(&u.Status, validation.Required),
		validation.Field(&u.CreatedAt, validation.Required),
		validation.Field(&u.UpdatedAt, validation.Required, validation.Min(u.CreatedAt)),
	)
}

func (u *User) ValidateUserForGoogleAuth() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.FamilyName, validation.Required, validation.RuneLength(2, 255)),
		validation.Field(&u.Name, validation.Required, validation.RuneLength(2, 255)),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.AvatarUrl, is.URL),
		validation.Field(&u.CapwallUrl, is.URL),
		validation.Field(&u.Privacy, validation.In(consts.PrivacyLevels...)),
		validation.Field(&u.Biography, validation.RuneLength(0, 500)),
		validation.Field(&u.AuthType, validation.In(consts.AuthTypes...)),
		validation.Field(&u.PostCount, validation.Min(0)),
		validation.Field(&u.FriendCount, validation.Min(0)),
		validation.Field(&u.Status, validation.Required),
		validation.Field(&u.CreatedAt, validation.Required),
		validation.Field(&u.UpdatedAt, validation.Required, validation.Min(u.CreatedAt)),
	)
}

func (u *UserUpdate) ValidateUserUpdate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.FamilyName, validation.RuneLength(2, 255)),
		validation.Field(&u.Name, validation.RuneLength(2, 255)),
		validation.Field(&u.PhoneNumber, validation.RuneLength(10, 14), validation.Match((regexp.MustCompile((`^\d+$`))))),
		validation.Field(&u.AvatarUrl, is.URL),
		validation.Field(&u.CapwallUrl, is.URL),
		validation.Field(&u.Privacy, validation.In(consts.PrivacyLevels...)),
		validation.Field(&u.Biography, validation.RuneLength(0, 500)),
		validation.Field(&u.AuthType, validation.In(consts.AuthTypes...)),
		validation.Field(&u.PostCount, validation.Min(0)),
		validation.Field(&u.FriendCount, validation.Min(0)),
	)
}

func NewUserLocal(
	familyName string,
	name string,
	email string,
	password string,
	phoneNumber string,
	birthday time.Time,
) (*User, error) {
	user := &User{
		ID:          uuid.New(),
		FamilyName:  familyName,
		Name:        name,
		Email:       email,
		Password:    &password,
		PhoneNumber: &phoneNumber,
		Birthday:    &birthday,
		AvatarUrl:   consts.AVATAR_URL,
		CapwallUrl:  consts.CAPWALL_URL,
		Privacy:     consts.PUBLIC,
		AuthType:    consts.LOCAL_AUTH,
		PostCount:   0,
		FriendCount: 0,
		Status:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := user.ValidateUser(); err != nil {
		return nil, err
	}

	return user, nil
}

func NewUserGoogle(
	familyName string,
	name string,
	email string,
	authGoogleId string,
	avatarUrl string,
) (*User, error) {
	user := &User{
		ID:           uuid.New(),
		FamilyName:   familyName,
		Name:         name,
		Email:        email,
		AuthGoogleId: &authGoogleId,
		AvatarUrl:    avatarUrl,
		CapwallUrl:   consts.CAPWALL_URL,
		Privacy:      consts.PUBLIC,
		AuthType:     consts.GOOGLE_AUTH,
		Birthday:     pointer.Ptr(time.Now()),
		PostCount:    0,
		FriendCount:  0,
		Status:       true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := user.ValidateUserForGoogleAuth(); err != nil {
		return nil, err
	}

	return user, nil
}
