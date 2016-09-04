package core

import (
	"math"
	"reflect"
	"testing"

	"github.com/yingtu/wukong/types"
	"github.com/yingtu/wukong/utils"
)

type DummyScoringFields struct {
	label   string
	counter int
	amount  float64
}

type DummyScoringCriteria struct {
	Threshold float64
}

func (criteria DummyScoringCriteria) Score(
	queryFields interface{}, doc types.IndexedDocument, fields interface{}) []float32 {
	if reflect.TypeOf(fields) == reflect.TypeOf(DummyScoringFields{}) {
		dsf := fields.(DummyScoringFields)
		value := float64(dsf.counter) + dsf.amount
		if reflect.TypeOf(queryFields) == reflect.TypeOf(DummyScoringFields{}) {
			qf := queryFields.(DummyScoringFields)
			value = math.Abs(value - float64(qf.counter) - qf.amount)
		}
		if criteria.Threshold != 0 && value > criteria.Threshold {
			return []float32{}
		}
		return []float32{float32(value)}
	}
	return []float32{}
}

func TestRankDocument(t *testing.T) {
	var ranker Ranker
	ranker.Init()
	ranker.AddDoc(1, DummyScoringFields{})
	ranker.AddDoc(3, DummyScoringFields{})
	ranker.AddDoc(4, DummyScoringFields{})

	scoredDocs, _ := ranker.Rank([]types.IndexedDocument{
		types.IndexedDocument{DocId: 1, BM25: 6},
		types.IndexedDocument{DocId: 3, BM25: 24},
		types.IndexedDocument{DocId: 4, BM25: 18},
	}, types.RankOptions{ScoringCriteria: types.RankByBM25{}}, false)
	utils.Expect(t, "[3 [24000 ]] [4 [18000 ]] [1 [6000 ]] ", scoredDocsToString(scoredDocs))

	scoredDocs, _ = ranker.Rank([]types.IndexedDocument{
		types.IndexedDocument{DocId: 1, BM25: 6},
		types.IndexedDocument{DocId: 3, BM25: 24},
		types.IndexedDocument{DocId: 2, BM25: 0},
		types.IndexedDocument{DocId: 4, BM25: 18},
	}, types.RankOptions{ScoringCriteria: types.RankByBM25{}, ReverseOrder: true}, false)
	// doc2因为没有AddDoc所以没有添加进来
	utils.Expect(t, "[1 [6000 ]] [4 [18000 ]] [3 [24000 ]] ", scoredDocsToString(scoredDocs))
}

func TestRankWithCriteria(t *testing.T) {
	var ranker Ranker
	ranker.Init()
	ranker.AddDoc(1, DummyScoringFields{
		label:   "label3",
		counter: 3,
		amount:  22.3,
	})
	ranker.AddDoc(2, DummyScoringFields{
		label:   "label4",
		counter: 1,
		amount:  2,
	})
	ranker.AddDoc(3, DummyScoringFields{
		label:   "label1",
		counter: 7,
		amount:  10.3,
	})
	ranker.AddDoc(4, DummyScoringFields{
		label:   "label1",
		counter: -1,
		amount:  2.3,
	})

	criteria := DummyScoringCriteria{}
	queryFields := DummyScoringFields{label: "ganga", counter: 5, amount: 10.3}
	scoredDocs, _ := ranker.Rank([]types.IndexedDocument{
		types.IndexedDocument{DocId: 1, TokenProximity: 6},
		types.IndexedDocument{DocId: 2, TokenProximity: -1},
		types.IndexedDocument{DocId: 3, TokenProximity: 24},
		types.IndexedDocument{DocId: 4, TokenProximity: 18},
	}, types.RankOptions{ScoringCriteria: criteria, QueryFields: queryFields, ReverseOrder: true}, false)
	utils.Expect(t, "[3 [2000 ]] [1 [10000 ]] [2 [12300 ]] [4 [14000 ]] ", scoredDocsToString(scoredDocs))

	criteria.Threshold = 10.0
	queryFields.amount = 1.3
	scoredDocs, _ = ranker.Rank([]types.IndexedDocument{
		types.IndexedDocument{DocId: 1, TokenProximity: 6},
		types.IndexedDocument{DocId: 2, TokenProximity: -1},
		types.IndexedDocument{DocId: 3, TokenProximity: 24},
		types.IndexedDocument{DocId: 4, TokenProximity: 18},
	}, types.RankOptions{ScoringCriteria: criteria, QueryFields: queryFields, ReverseOrder: true}, false)
	utils.Expect(t, "[2 [3300 ]] [4 [5000 ]] ", scoredDocsToString(scoredDocs))
}

func TestRemoveDoc(t *testing.T) {
	var ranker Ranker
	ranker.Init()
	ranker.AddDoc(1, DummyScoringFields{
		label:   "label3",
		counter: 3,
		amount:  22.3,
	})
	ranker.AddDoc(2, DummyScoringFields{
		label:   "label4",
		counter: 1,
		amount:  2,
	})
	ranker.AddDoc(3, DummyScoringFields{
		label:   "label1",
		counter: 7,
		amount:  10.3,
	})
	ranker.RemoveDoc(3)

	criteria := DummyScoringCriteria{}
	scoredDocs, _ := ranker.Rank([]types.IndexedDocument{
		types.IndexedDocument{DocId: 1, TokenProximity: 6},
		types.IndexedDocument{DocId: 2, TokenProximity: -1},
		types.IndexedDocument{DocId: 3, TokenProximity: 24},
		types.IndexedDocument{DocId: 4, TokenProximity: 18},
	}, types.RankOptions{ScoringCriteria: criteria}, false)
	utils.Expect(t, "[1 [25300 ]] [2 [3000 ]] ", scoredDocsToString(scoredDocs))
}
