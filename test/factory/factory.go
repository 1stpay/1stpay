package factory

import "gorm.io/gorm"

type TestFactory struct {
	db          *gorm.DB
	UserFactory *TestUserFactory
}

func NewTestFactory(db *gorm.DB) *TestFactory {
	return &TestFactory{
		db:          db,
		UserFactory: NewTestUserFactory(db),
	}
}
