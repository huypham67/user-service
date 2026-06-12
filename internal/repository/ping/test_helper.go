package ping

import (
	"testing"

	"github.com/huypham67/bookmark-common/pkg/sqldb"
	"gorm.io/gorm"
)

func newTestSQLDBPinger(t *testing.T) (Pinger, *gorm.DB) {
	t.Helper()

	mockDB := sqldb.NewMock(t)
	pinger := NewSQLDB(mockDB)

	return pinger, mockDB
}
