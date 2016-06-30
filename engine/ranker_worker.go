package engine

import (
	"github.com/yingtu/wukong/types"
)

type rankerAddDocRequest struct {
	docId  uint64
	fields interface{}
}

type rankerRemoveDocRequest struct {
	docId uint64
}

type rankerLookupDocRequest struct {
	docId               uint64
	rankerReturnChannel chan interface{}
}

type rankerRankRequest struct {
	docs                []types.IndexedDocument
	options             types.RankOptions
	rankerReturnChannel chan rankerReturnRequest
	countDocsOnly       bool
}

type rankerReturnRequest struct {
	docs    types.ScoredDocuments
	numDocs int
}

func (engine *Engine) rankerAddDocWorker(shard int) {
	for {
		request := <-engine.rankerAddDocChannels[shard]
		engine.rankers[shard].AddDoc(request.docId, request.fields)
	}
}

func (engine *Engine) rankerRemoveDocWorker(shard int) {
	for {
		request := <-engine.rankerRemoveDocChannels[shard]
		engine.rankers[shard].RemoveDoc(request.docId)
	}
}

func (engine *Engine) rankerUpdateDocWorker(shard int) {
	for {
		request := <-engine.rankerUpdateDocChannels[shard]
		engine.rankers[shard].UpdateDoc(request.docId, request.fields)
	}
}

func (engine *Engine) rankerLookupDocWorker(shard int) {
	for {
		request := <-engine.rankerLookupDocChannels[shard]
		fields := engine.rankers[shard].LookupDoc(request.docId)
		request.rankerReturnChannel <- fields
	}
}

func (engine *Engine) rankerRankWorker(shard int) {
	for {
		request := <-engine.rankerRankChannels[shard]
		if request.options.MaxOutputs != 0 {
			request.options.MaxOutputs += request.options.OutputOffset
		}
		request.options.OutputOffset = 0
		outputDocs, numDocs := engine.rankers[shard].Rank(request.docs, request.options, request.countDocsOnly)
		request.rankerReturnChannel <- rankerReturnRequest{docs: outputDocs, numDocs: numDocs}
	}
}
