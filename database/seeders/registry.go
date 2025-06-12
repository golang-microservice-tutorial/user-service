package seeders

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

type SeederRegistry interface {
	Run()
}

func NewSeederRegistry(db *gorm.DB) *Registry {
	return &Registry{db: db}
}

func (r *Registry) Run() {
	RunRoleSeeder(r.db)
	RunUserSeeder(r.db)
	logrus.Info("Database seeding completed successfully")
}
