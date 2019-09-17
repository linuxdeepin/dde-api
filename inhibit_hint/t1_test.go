package inhibit_hint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSetName(t *testing.T) {
	o := New("domain")
	o.SetName("name1")
	assert.Equal(t, "name1", o.getName("why1"))
	assert.Equal(t, "name1", o.getName("why2"))

	fn := func(why string) string {
		switch why {
		case "why1":
			return "name1"
		case "why2":
			return "name2"
		}
		return ""
	}
	o.SetNameFunc(fn)
	assert.Equal(t, "name1", o.getName("why1"))
	assert.Equal(t, "name2", o.getName("why2"))
}

func TestGetSetIcon(t *testing.T) {
	o := New("domain")
	o.SetIcon("icon1")
	assert.Equal(t, "icon1", o.getIcon("why1"))
	assert.Equal(t, "icon1", o.getIcon("why2"))

	fn := func(why string) string {
		switch why {
		case "why1":
			return "icon1"
		case "why2":
			return "icon2"
		}
		return ""
	}
	o.SetIconFunc(fn)
	assert.Equal(t, "icon1", o.getIcon("why1"))
	assert.Equal(t, "icon2", o.getIcon("why2"))
}
