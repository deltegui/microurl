package persistence_test

import (
	"microurl/internal"
	"microurl/internal/persistence"
	"microurl/testutils"
	"testing"
)

func TestShouldInsertURLEntries(t *testing.T) {
	type data struct {
		name string
		url  internal.URL
	}
	runTest := func(current data) {
		t.Run(current.name, func(t *testing.T) {
			testutils.DBTransaction(func(conn persistence.Connection) {
				repo := persistence.NewGormURLRepository(conn)
				if err := repo.Save(&current.url); err != nil {
					t.Error(err)
				}
				url, err := repo.FindByID(int(current.url.ID))
				if err != nil {
					t.Error(err)
				}
				if url != current.url {
					t.Error("Expected url ", current.url, " to be equal to ", url)
				}
			})
		})
	}
	cases := []data{}
	for _, c := range cases {
		runTest(c)
	}
}
