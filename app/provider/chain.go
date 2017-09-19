package provider

import (
	"github.com/MoonBabyLabs/boom/app/service/chain"
)

type ChainProvider struct {

}


func (c ChainProvider) Construct() chain.BoomChainHandler {
	bc := chain.BoomChain{}.SetBlock(chain.BoomBlock{})
	bc.SetBlocks(make([]map[string][]byte, 0))

	return bc
}