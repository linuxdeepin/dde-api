package drandr

import (
	"testing"

	x "github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/ext/randr"
	"github.com/stretchr/testify/assert"
)

var outputInfo1 = OutputInfo{
	Name:     "OutputInfo1",
	Id:       1,
	MmWidth:  1920,
	MmHeight: 1080,
	Crtc: CrtcInfo{
		Id:        1,
		Mode:      1,
		X:         0,
		Y:         0,
		Width:     1920,
		Height:    1080,
		Rotation:  60,
		Reflect:   1,
		Rotations: nil,
		Reflects:  nil,
	},
	Connection:     true,
	Timestamp:      0,
	EDID:           nil,
	Clones:         nil,
	Crtcs:          nil,
	Modes:          nil,
	PreferredModes: nil,
}
var outputInfo2 = OutputInfo{
	Name:     "OutputInfo2",
	Id:       2,
	MmWidth:  1920,
	MmHeight: 1080,
	Crtc: CrtcInfo{
		Id:        2,
		Mode:      1,
		X:         0,
		Y:         0,
		Width:     1920,
		Height:    1080,
		Rotation:  60,
		Reflect:   1,
		Rotations: nil,
		Reflects:  nil,
	},
	Connection:     false,
	Timestamp:      0,
	EDID:           nil,
	Clones:         nil,
	Crtcs:          nil,
	Modes:          nil,
	PreferredModes: nil,
}
var outputInfo3 = OutputInfo{
	Name:     "OutputInfo3",
	Id:       3,
	MmWidth:  1920,
	MmHeight: 1080,
	Crtc: CrtcInfo{
		Id:        3,
		Mode:      1,
		X:         0,
		Y:         0,
		Width:     1920,
		Height:    1080,
		Rotation:  60,
		Reflect:   1,
		Rotations: nil,
		Reflects:  nil,
	},
	Connection:     true,
	Timestamp:      0,
	EDID:           nil,
	Clones:         nil,
	Crtcs:          nil,
	Modes:          nil,
	PreferredModes: nil,
}
var outputInfos = OutputInfos{
	outputInfo1,
	outputInfo2,
	outputInfo3,
}

func Test_OutputInfosQuery(t *testing.T) {

	assert.Equal(t, outputInfo1, outputInfos.Query(1))
	assert.Equal(t, outputInfo2, outputInfos.Query(2))
	assert.Equal(t, outputInfo3, outputInfos.Query(3))
	assert.Equal(t, OutputInfo{}, outputInfos.Query(4))

}

func Test_QueryByName(t *testing.T) {

	assert.Equal(t, outputInfo1, outputInfos.QueryByName("OutputInfo1"))
	assert.Equal(t, outputInfo2, outputInfos.QueryByName("OutputInfo2"))
	assert.Equal(t, outputInfo3, outputInfos.QueryByName("OutputInfo3"))
	assert.Equal(t, OutputInfo{}, outputInfos.QueryByName("OutputInfo4"))

}

func Test_ListNames(t *testing.T) {

	assert.Equal(t, []string{
		"OutputInfo1",
		"OutputInfo2",
		"OutputInfo3",
	}, outputInfos.ListNames())

}

func Test_ListConnectionOutputs(t *testing.T) {

	assert.Equal(t, outputInfo1, outputInfos.ListConnectionOutputs()[0])
	assert.Equal(t, outputInfo3, outputInfos.ListConnectionOutputs()[1])

}

func Test_outputsToRandrIdList(t *testing.T) {

	list := outputsToRandrIdList([]randr.Output{1, 3, 5, 6, 9})
	assert.Equal(t, randrIdList{1, 3, 5, 6, 9}, list)

}

func Test_crtcsToRandrIdList(t *testing.T) {

	list := crtcsToRandrIdList([]randr.Crtc{1, 3, 5, 6, 9})
	assert.Equal(t, randrIdList{1, 3, 5, 6, 9}, list)

}

func Test_modesToRandrIdList(t *testing.T) {

	list := modesToRandrIdList([]randr.Mode{1, 3, 5, 6, 9})
	assert.Equal(t, randrIdList{1, 3, 5, 6, 9}, list)

}

func Test_getOutputEDID(t *testing.T) {
	xConn, err := x.NewConn()
	if err != nil {
		t.Skip(err)
	}
	root := xConn.GetDefaultScreen().Root
	resource, err := randr.GetScreenResources(xConn, root).Reply(xConn)
	if err != nil {
		t.Skip(err)
	}
	for _, output := range resource.Outputs {
		t.Run("Test_getOutputEDID", func(t *testing.T) {
			_, err := getOutputEDID(xConn, output)
			assert.NoError(t, err)
		})
	}

}
