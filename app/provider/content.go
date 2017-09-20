package provider

import (
	"github.com/MoonBabyLabs/boom/app/service/content"
)

type Content struct {

}

func (c Content) Construct() content.Manager {
	cnt := content.Default{}.
		SetDatastore(Db{}.Construct()).
		SetChain(ChainProvider{}.Construct()).
		SetFileManager(Filemanager{}.Construct()).
		SetUuidGenerator(Uuid{}.Construct())

	return cnt
}
