package user

import (
	"testing"

	"github.com/huypham67/user-service/internal/test/fixtures"
	"gorm.io/gorm"
)

func newTestRepository(t *testing.T) (Repository, *gorm.DB) {
	t.Helper()

	testDB := fixtures.NewTestDB(t, &fixtures.UserTestDB{})
	repo := NewRepository(testDB)

	return repo, testDB
}
