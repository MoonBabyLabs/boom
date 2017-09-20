package provider

import (
	"github.com/MoonBabyLabs/boom/app/service/chain"
)

type ChainProvider struct {

}


func (c ChainProvider) Construct() chain.BoomChainHandler {
	bc := chain.BoomChain{}.SetBlock(chain.BoomBlock{}).SetBlocks(make([]map[string]map[string]interface{}, 0))

	return bc
}