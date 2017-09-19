package content_test

import (
	"testing"
	"github.com/MoonBabyLabs/boom/app/service/content"
	"github.com/MoonBabyLabs/boom/tests/mocks"
)
type defaultTest struct {
	*testing.T
}

func TestSetDatastore(t *testing.T) {
	defaultContent := content.Default{}
	d := mocks.Datastore{}
	e := d.Init("a", "b")
	s := make(map[string]interface{})
	s["sample"] = "content"

	b := defaultContent.SetDatastore(e).Datastore().Find("a", "b")

	if s["sample"] != b["sample"] {
		t.Errorf("could not get datastore content")
	}
}

func BenchmarkSetdatastore(t *testing.B) {
	dc := content.Default{}
	d := mocks.Datastore{}
	e := d.Init("a", "b")
	dc.SetDatastore(e).Datastore()
}