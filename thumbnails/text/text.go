// Text thumbnail generator
package text

import (
	"fmt"
	. "pkg.deepin.io/dde/api/thumbnails/loader"
	"pkg.deepin.io/lib/mime"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	TextTypeText = "text/plain"
	TextTypeC    = "text/x-c"
	TextTypeCpp  = "text/x-cpp"
	TextTypeGo   = "text/x-go"
)

func init() {
	for _, ty := range SupportedTypes() {
		Register(ty, genTextThumbnail)
	}
}

func SupportedTypes() []string {
	return []string{
		TextTypeText,
		TextTypeC,
		TextTypeCpp,
		TextTypeGo,
	}
}

func GenThumbnail(src string, width, height int) (string, error) {
	ty, err := mime.Query(src)
	if err != nil {
		return "", err
	}

	if !IsStrInList(ty, SupportedTypes()) {
		return "", fmt.Errorf("Not supported type: %v", ty)
	}

	return genTextThumbnail(src, "", width, height)
}

func genTextThumbnail(src, bg string, width, height int) (string, error) {
	dest, err := GetThumbnailDest(src, width, height)
	if err != nil {
		return "", err
	}
	if dutils.IsFileExist(dest) {
		return dest, nil
	}

	return doGenThumbnail(dutils.DecodeURI(src), dest, width, height)
}
