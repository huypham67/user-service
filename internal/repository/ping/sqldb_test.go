package ping

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSqlDBPinger_Ping(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		verify func(*testing.T, error)
	}{
		{
			name: "should successfully ping the database",
			verify: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "should return error when the database is closed",
			verify: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			pinger, mockDB := newTestSQLDBPinger(t)

			if tc.name == "should return error when the database is closed" {
				sqlDB, err := mockDB.DB()
				require.NoError(t, err)
				require.NoError(t, sqlDB.Close())
			}

			err := pinger.Ping(ctx)

			tc.verify(t, err)
		})
	}
}
