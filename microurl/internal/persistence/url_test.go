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
					url, err := populator.URLRepo.FindByID(expected.ID)
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
			name: "should delete one element",
			exptectedMissing: []internal.URL{
				{
					ID:       2,
					Original: "https://hello.com/hola",
					Owner:    "manolo",
				},
			},
		},
		{
			name: "should delete many elements",
			exptectedMissing: []internal.URL{
				{
					ID:       1,
					Original: "http://youtube.com/hola",
					Owner:    "manolo",
				},
				{
					ID:       3,
					Original: "https://youtube.com/xasfhasd",
					Owner:    "paola",
				},
				{
					ID:       4,
					Original: "https://manolo.com/manolo",
					Owner:    "ambrosio",
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
					_, err := populator.URLRepo.FindByID(expected.ID)
					if err == nil {
						t.Error("Expected url", expected, "to be deleted")
						return
					}
				}
			})
		})
	}
}

func TestIfTheElementAlreadyExistsShouldRewriteIt(t *testing.T) {
	url := internal.URL{
		ID:       0,
		Original: "https://hello.com/hola",
		Owner:    "manolo",
		Times:    0,
	}
	expected := internal.URL{
		ID:       1,
		Original: "https://hello.com/hola",
		Owner:    "manolo",
		Times:    2,
	}
	testutils.DBTransaction(func(populator testutils.Populator) {
		populator.PopulateUsers()
		if err := populator.URLRepo.Save(&url); err != nil {
			t.Error(err)
			return
		}
		url.Times += 2
		if err := populator.URLRepo.Save(&url); err != nil {
			t.Error(err)
			return
		}
		result, err := populator.URLRepo.FindByID(1)
		if err != nil {
			t.Error(err)
			return
		}
		if result != expected {
			t.Error("Expected", expected, "but got", result)
		}
	})
}

func TestGetAll(t *testing.T) {
	expected := []internal.URL{
		{
			ID:       1,
			Original: "http://youtube.com/hola",
			Owner:    "manolo",
			Times:    0,
		},
		{
			ID:       2,
			Original: "https://hello.com/hola",
			Owner:    "manolo",
			Times:    0,
		},
	}
	testutils.DBTransaction(func(populator testutils.Populator) {
		populator.PopulateAll()
		urls := populator.URLRepo.GetAllForUser("manolo")
		if len(urls) != len(expected) {
			t.Error("Expected length of", len(expected), "but have", len(urls))
			return
		}
		for i, url := range urls {
			if url != expected[i] {
				t.Errorf("[For index %d] Expected %+v to be equal to %+v", i, expected[i], url)
			}
		}
	})
}
