package provider

import (
	"github.com/MoonBabyLabs/boom/app/service/content"
	"log"
)

type Content struct {

}

func (c Content) Construct() content.Manager {
	cnt := content.Default{}.SetDatastore(Db{}.Construct()).SetChain(ChainProvider{}.Construct()).SetFileManager(Filemanager{}.Construct())
	log.Print(cnt.Chain().Block())
	return cnt
}
