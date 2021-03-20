package sqlparser

import (
	"os"
	"testing"

	"github.com/dvasilas/proteus/internal/libqpu"
	"github.com/dvasilas/proteus/internal/proto/qpu"
	"github.com/dvasilas/proteus/internal/proto/qpuapi"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	returnCode := m.Run()
	os.Exit(returnCode)
}

var filterTests = []struct {
	querySQL       string
	expectedQueryI *qpuapi.ASTQuery
}{
	{
		"select * from t where x = 42",
		&qpuapi.ASTQuery{
			Projection: []string{"*"},
			Table:      "t",
			Predicate: []*qpu.AttributePredicate{
				&qpu.AttributePredicate{
					Attr:   libqpu.Attribute("x", nil),
					Type:   qpu.AttributePredicate_EQ,
					Lbound: libqpu.ValueInt(42),
					Ubound: libqpu.ValueInt(42),
				},
			},
			TsPredicate: libqpu.SnapshotTimePredicate(
				libqpu.SnapshotTime(qpu.SnapshotTime_LATEST, nil, true),
				libqpu.SnapshotTime(qpu.SnapshotTime_LATEST, nil, true),
			),
		},
	},
	{
		"SELECT title, description, short_id, user_id, vote_sum FROM qpu ORDER BY vote_sum DESC LIMIT 5",
		&qpuapi.ASTQuery{
			Projection: []string{"title", "description", "short_id", "user_id", "vote_sum"},
			Table:      "qpu",
			OrderBy: &qpuapi.OrderBy{
				AttributeName: "vote_sum",
				Direction:     qpuapi.OrderBy_DESC,
			},
			Limit: int64(5),
			TsPredicate: libqpu.SnapshotTimePredicate(
				libqpu.SnapshotTime(qpu.SnapshotTime_LATEST, nil, true),
				libqpu.SnapshotTime(qpu.SnapshotTime_LATEST, nil, true),
			),
		},
	},
}

func TestParse(t *testing.T) {
	for _, tt := range filterTests {
		queryI, _ := Parse(tt.querySQL)
		assert.Equal(t, tt.expectedQueryI, queryI.Q, "")
	}
}
