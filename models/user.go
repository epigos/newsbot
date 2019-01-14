package models

import (
	"time"

	"cloud.google.com/go/datastore"
)

// UserKind kind name for facebook user
const UserKind = "Users"

// UserActionKind kind name for user actions
const UserActionKind = "UserActions"

// User represent facebook user
type User struct {
	ID        string    `json:"id" datastore:"-"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Avatar    string    `json:"profile_pic" datastore:",noindex"`
	Locale    string    `json:"locale" datastore:",noindex"`
	TimeZone  int32     `json:"timezone" datastore:",noindex"`
	Gender    string    `json:"gender"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
}

// UserAction record user actions
type UserAction struct {
	ID      string         `json:"id" datastore:"-"`
	UserKey *datastore.Key `json:"bot_user_id"`
	ItemKey *datastore.Key `json:"item_id"`
	Action  string         `json:"action"`
	Created time.Time      `json:"created"`
	Updated time.Time      `json:"updated"`
}

// Key get key for article
func (m *User) Key() *datastore.Key {
	// if there is no Id, we want to generate an "incomplete"
	// one and let datastore determine the key/Id for us
	if m.ID == "" {
		return DS.NewKey(UserKind)
	}

	// if Id is already set, we'll just build the Key based
	// on the one provided.
	return datastore.NameKey(UserKind, m.ID, nil)
}

// SetID set id
func (m *User) SetID(key *datastore.Key) {
	m.ID = key.Name
}

// NewUser returns new model
func NewUser(id, firstName, lastName, avatar, locale, gender string, timezone int32) *User {
	return &User{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Avatar:    avatar,
		Locale:    locale,
		TimeZone:  timezone,
		Gender:    gender,
		Created:   time.Now(),
	}
}

func (m *User) String() string {
	return m.FirstName + " " + m.LastName
}

// GetUser article
func GetUser(id string) (*User, error) {
	entity := User{ID: id}
	err := DS.GetByKey(&entity)
	return &entity, err
}

// Save users
func (m *User) Save() {
	DS.Logger.Info("Saving user:", m)
	DS.Save(m)
}

// GetUserKey get user key
func GetUserKey(id string) *datastore.Key {
	entity := User{ID: id}
	return entity.Key()
}

// Key get key for article
func (m *UserAction) Key() *datastore.Key {
	// if there is no Id, we want to generate an "incomplete"
	// one and let datastore determine the key/Id for us
	if m.ID == "" {
		return DS.NewKey(UserActionKind)
	}

	// if Id is already set, we'll just build the Key based
	// on the one provided.
	key, err := datastore.DecodeKey(m.ID)
	if err != nil {
		DS.Logger.Error("Key not found:", err)
	}
	return key
}

// SetID set id
func (m *UserAction) SetID(key *datastore.Key) {
	m.ID = key.Encode()
}

// NewUserAction creates a new bot user action
func NewUserAction(userID string, itemKey *datastore.Key, action string) *UserAction {
	return &UserAction{
		UserKey: GetUserKey(userID),
		ItemKey: itemKey,
		Action:  action,
	}
}

// Save UserAction
func (m *UserAction) Save() {
	DS.Logger.Info("Saving bot user action:", m)
	DS.Save(m)
}
