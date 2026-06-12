package fixtures

import (
	"github.com/huypham67/user-service/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const TestPassword = "password123"

// UserTestDB is a test database struct for user-related tests. It embeds baseTestDB to inherit common database setup and teardown functionalities.
type UserTestDB struct {
	baseTestDB
}

// NewUserTestDB creates a new UserTestDB instance with the given gorm.DB connection.
func (u *UserTestDB) MigrateDB() error {
	return u.db.AutoMigrate(&model.User{})
}

// SeedData populates the database with predefined test users. It uses a session with SkipHooks to bypass any model hooks during seeding.
func (u *UserTestDB) SeedData() error {
	db := u.db.Session(&gorm.Session{SkipHooks: true})

	users := []model.User{
		{
			BaseModel: model.BaseModel{
				ID: "user-uuid-1",
			},
			DisplayName: "Test User 1",
			Username:    "testuser1",
			Email:       "testuser1@gmail.com",
			Password:    hashPassword(TestPassword),
		},
		{
			BaseModel: model.BaseModel{
				ID: "user-uuid-2",
			},
			DisplayName: "Test User 2",
			Username:    "testuser2",
			Email:       "testuser2@gmail.com",
			Password:    hashPassword(TestPassword),
		},
		{
			BaseModel: model.BaseModel{
				ID: "user-uuid-3",
			},
			DisplayName: "Test User 3",
			Username:    "testuser3",
			Email:       "testuser3@gmail.com",
			Password:    hashPassword(TestPassword),
		},
	}

	err := db.CreateInBatches(users, 10).Error
	if err != nil {
		return err
	}

	return nil
}

func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		panic(err)
	}

	return string(hash)
}
