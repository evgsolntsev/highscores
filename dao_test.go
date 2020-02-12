package main

import (
	"context"
	"testing"

	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDAO(t *testing.T) {
	scores := []Score{{
		Name:  "Docker",
		Score: 70,
	}, {
		Name:  "Docker",
		Score: 80,
	}, {
		Name:  "Evgsol",
		Score: 50,
	}}

	session, err := mgo.Dial("mongodb://localhost:27017")
	require.Nil(t, err)

	ctx := context.Background()
	dao := NewDAO(ctx, session)
	defer dao.RemoveAll(ctx)

	for _, s := range scores {
		err = dao.Insert(ctx, s)
		require.Nil(t, err)
	}

	type testcase struct {
		q        int
		expected []Score
	}

	for _, test := range []testcase{{
		q: 1,
		expected: []Score{{Name: "Docker", Score: 80}},
	}, {
		q: 2,
		expected: []Score{
			{Name: "Docker", Score: 80},
			{Name: "Docker", Score: 70},
		},
	}, {
		q: 3,
		expected: []Score{
			{Name: "Docker", Score: 80},
			{Name: "Docker", Score: 70},
			{Name: "Evgsol", Score: 50},
		},
	}, {
		q: 4,
		expected: []Score{
			{Name: "Docker", Score: 80},
			{Name: "Docker", Score: 70},
			{Name: "Evgsol", Score: 50},
		},
	}} {
		realScores, err := dao.GetTop(ctx, test.q)
		assert.NoError(t, err)
		if assert.Len(t, realScores, len(test.expected)) {
			for i, realScore := range realScores {
				assert.Equal(t, realScore, test.expected[i])
			}
		}
	}
}
