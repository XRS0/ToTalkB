package codes

import (
	"fmt"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

func Generate(contains string, path string) {
	qrc, err := qrcode.New(contains)
	if err != nil {
		fmt.Printf("could not generate QRCode: %v", err)
		return
	}

	w, err := standard.New(path)
	if err != nil {
		fmt.Printf("standard.New failed: %v", err)
		return
	}

	if err = qrc.Save(w); err != nil {
		fmt.Printf("could not save image: %v", err)
	}
}
