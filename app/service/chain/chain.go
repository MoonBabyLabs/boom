package chain

import (
	"log"
	"encoding/json"
)

type BoomChain struct {
	blocks []map[string]map[string]interface{}
	block BlockHandler
}

type BoomChainHandler interface {
	Init(blockHandler BlockHandler) BoomChainHandler
	New(genesis BlockHandler) BoomChainHandler
	AddBlock(block BlockHandler) BoomChainHandler
	Blocks() []map[string]map[string]interface{}
	//Version() BoomChain  @todo Will be implemented in future version
	//NewVersion() BoomChain @todo will be implemented in future version
	Latest() string
	Block() BlockHandler
	SetBlock(block BlockHandler) BoomChainHandler
	SetBlocks(blocks []map[string]map[string]interface{}) BoomChainHandler
}

func (bc BoomChain) Init(handler BlockHandler) BoomChainHandler {
	blocks := make([]map[string]map[string]interface{}, 0)
	bc.SetBlocks(blocks)
	return bc.SetBlock(handler)
}

func (t BoomChain) Block() BlockHandler {
	return t.block
}

func (t BoomChain) SetBlock(block BlockHandler) BoomChainHandler {
	log.Print(block)
	t.block = block

	return t
}

// New instantiates a new chain
func (t BoomChain) New(genesis BlockHandler) BoomChainHandler {
	t.AddBlock(genesis)

	return t
}

func (b BoomChain) Latest() string {
	block := b.blocks[len(b.blocks) -1]
	for k, _ := range block {
		return k
	}

	return ""
}

func (t BoomChain) AddBlock(block BlockHandler) BoomChainHandler {
	item := make(map[string]map[string]interface{})
	data := make(map[string]interface{})
	json.Unmarshal(block.Data(), &data)
	item[block.HashString()] = data
	t.blocks = append(t.blocks, item)

	return t
}

func (t BoomChain) Blocks() []map[string]map[string]interface{} {
	return t.blocks
}

func (t BoomChain) SetBlocks(blocks []map[string]map[string]interface{}) BoomChainHandler {
	t.blocks = blocks

	return t
}

