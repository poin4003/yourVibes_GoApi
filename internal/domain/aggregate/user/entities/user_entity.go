package entities

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"time"
)

type User struct {
	ID           uuid.UUID           `validate:"required,uuid4"`
	FamilyName   string              `validate:"required,min=2"`
	Name         string              `validate:"required,min=2"`
	Email        string              `validate:"required,email"`
	Password     string              `validate:"required,min=8"`
	PhoneNumber  string              `validate:"required,min=10,max=14,numeric"`
	Birthday     time.Time           `validate:"required,date"`
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
	Setting      *Setting            `validate:"required"`
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
	setting *Setting,
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
		Setting:     setting,
	}
	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUpdate) SetUpdatedAt() {
	now := time.Now()
	u.UpdatedAt = &now
}

func (u *UserUpdate) UpdateFamilyName(familyName *string) error {
	u.FamilyName = familyName
	u.SetUpdatedAt()
	return u.ValidateUserUpdate()
}

func (u *UserUpdate) UpdateName(name *string) error {
	u.Name = name
	u.SetUpdatedAt()
	return u.ValidateUserUpdate()
}

func (u *UserUpdate) UpdateEmail(email *string) error {
	u.Email = email
	u.SetUpdatedAt()
	return u.ValidateUserUpdate()
}

func (u *UserUpdate) UpdatePassword(password *string) error {
	u.Password = password
	u.SetUpdatedAt()
	return u.ValidateUserUpdate()
}

func (u *UserUpdate) UpdatePhoneNumber(phoneNumber *string) error {
	u.PhoneNumber = phoneNumber
	u.SetUpdatedAt()
	return u.ValidateUserUpdate()
}

func (u *UserUpdate) UpdateBirthday(birthday *time.Time) error {
	u.Birthday = birthday
	u.SetUpdatedAt()
	return u.ValidateUserUpdate()
}

func (u *UserUpdate) UpdateAvatarUrl(avatarUrl *string) error {
	u.AvatarUrl = avatarUrl
	u.SetUpdatedAt()
	return u.ValidateUserUpdate()
}

func (u *UserUpdate) UpdateCapwallUrl(capwallUrl *string) error {
	u.CapwallUrl = capwallUrl
	u.SetUpdatedAt()
	return u.ValidateUserUpdate()
}

func (u *UserUpdate) UpdatePrivacy(privacy *consts.PrivacyLevel) error {
	if !consts.IsValidPrivacyLevel(*privacy) {
		return errors.New("invalid privacy level")
	}
	u.Privacy = privacy
	u.SetUpdatedAt()
	return u.ValidateUserUpdate()
}

func (u *UserUpdate) UpdateBiography(biography *string) error {
	u.Biography = biography
	u.SetUpdatedAt()
	return u.ValidateUserUpdate()
}

func (u *UserUpdate) UpdateAuthType(authType *consts.AuthType) error {
	u.AuthType = authType
	u.SetUpdatedAt()
	return u.ValidateUserUpdate()
}

func (u *UserUpdate) UpdateAuthGoogleId(authGoogleId *string) error {
	u.AuthGoogleId = authGoogleId
	u.SetUpdatedAt()
	return u.ValidateUserUpdate()
}

func (u *UserUpdate) IncreasePostCount() error {
	if u.PostCount == nil {
		initialCount := 0
		u.PostCount = &initialCount
	}
	*u.PostCount++
	u.SetUpdatedAt()
	return u.ValidateUserUpdate()
}

func (u *UserUpdate) DecreasePostCount() error {
	if u.PostCount == nil {
		initialCount := 0
		u.PostCount = &initialCount
	}
	*u.PostCount--
	u.SetUpdatedAt()
	return u.ValidateUserUpdate()
}

func (u *UserUpdate) IncreaseFriendCount() error {
	if u.FriendCount == nil {
		initialCount := 0
		u.FriendCount = &initialCount
	}
	*u.FriendCount++
	u.SetUpdatedAt()
	return u.ValidateUserUpdate()
}

func (u *UserUpdate) DecreaseFriendCount() error {
	if u.FriendCount == nil {
		initialCount := 0
		u.FriendCount = &initialCount
	}
	*u.FriendCount--
	u.SetUpdatedAt()
	return u.ValidateUserUpdate()
}

func (u *UserUpdate) Activate() error {
	*u.Status = true
	u.SetUpdatedAt()
	return u.ValidateUserUpdate()
}

func (u *UserUpdate) Deactivate() error {
	*u.Status = false
	u.SetUpdatedAt()
	return u.ValidateUserUpdate()
}
