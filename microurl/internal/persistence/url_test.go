package persistence_test

import (
	"microurl/internal"
	"microurl/testutils"
	"testing"
)

func TestShouldInsertURLEntries(t *testing.T) {
	type data struct {
		name       string
		url        []internal.URL
		insertions []internal.URL
	}
	cases := []data{
		{
			name: "should insert one element",
			url: []internal.URL{
				{
					ID:       1,
					Original: "https://hello.com/hola",
					Owner:    "manolo",
				},
			},
			insertions: []internal.URL{
				{
					ID:       0,
					Original: "https://hello.com/hola",
					Owner:    "manolo",
				},
			},
		},
		{
			name: "should insert many element",
			url: []internal.URL{
				{
					ID:       1,
					Original: "https://hello.com/hola",
					Owner:    "manolo",
				},
				{
					ID:       2,
					Original: "https://youtube.com/xasfhasd",
					Owner:    "paola",
				},
				{
					ID:       3,
					Original: "https://manolo.com/manolo",
					Owner:    "ambrosio",
				},
			},
			insertions: []internal.URL{
				{
					ID:       0,
					Original: "https://hello.com/hola",
					Owner:    "manolo",
				},
				{
					ID:       0,
					Original: "https://youtube.com/xasfhasd",
					Owner:    "paola",
				},
				{
					ID:       0,
					Original: "https://manolo.com/manolo",
					Owner:    "ambrosio",
				},
			},
		},
	}
	for _, current := range cases {
		t.Run(current.name, func(t *testing.T) {
			testutils.DBTransaction(func(populator testutils.Populator) {
				populator.PopulateUsers()
				for _, insert := range current.insertions {
					if err := populator.URLRepo.Save(&insert); err != nil {
						t.Error(err)
						return
					}
				}
				for _, expected := range current.url {
					url, err := populator.URLRepo.FindByID(int(expected.ID))
					if err != nil {
						t.Error(err)
						return
					}
					if url != expected {
						t.Error("Expected url", current.url, ", but have ", url)
					}
				}
			})
		})
	}
}

func TestShouldDeleteUrls(t *testing.T) {
	type data struct {
		name             string
		exptectedMissing []internal.URL
	}
	cases := []data{
		{
			name: "should delte one element",
			exptectedMissing: []internal.URL{
				{
					ID:       2,
					Original: "https://hello.com/hola",
					Owner:    "manolo",
				},
			},
		},
	}
	for _, current := range cases {
		t.Run(current.name, func(t *testing.T) {
			testutils.DBTransaction(func(populator testutils.Populator) {
				populator.PopulateAll()
				for _, url := range current.exptectedMissing {
					if err := populator.URLRepo.Delete(url); err != nil {
						t.Error(err)
						return

					}
				}
				for _, expected := range current.exptectedMissing {
					_, err := populator.URLRepo.FindByID(int(expected.ID))
					if err == nil {
						t.Error("Expected url", expected, "to be deleted")
						return
					}
				}
			})
		})
	}
}
