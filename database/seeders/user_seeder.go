package seeders

import (
	"user-service/constants"
	"user-service/domain/models"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RunUserSeeder(db *gorm.DB) {
	password, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		panic("Failed to hash password: " + err.Error())
	}

	user := models.User{
		UUID:        uuid.New(),
		Name:        "Admin",
		Username:    "admin",
		Password:    string(password),
		PhoneNumber: "08123456789",
		Email:       "admin@admin.com",
		RoleID:      constants.Admin,
	}

	err = db.FirstOrCreate(&user, models.User{Username: user.Username}).Error
	if err != nil {
		logrus.Errorf("Failed to seed user %s: %v", user.Username, err)
		panic(err)
	}
	logrus.Infof("User %s seeded successfully", user.Username)
}
