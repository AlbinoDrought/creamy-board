package thumbnailer

import (
	"bytes"
	"errors"
	"image"
	"io"

	"github.com/gabriel-vasile/mimetype"
	"github.com/nfnt/resize"
	"go.albinodrought.com/creamy-board/internal/log"

	_ "image/gif"
	"image/jpeg"
	"image/png"
)

type Thumb struct {
	io.ReadCloser
	Bytes int
	MIME  *mimetype.MIME
}

func supportedMime(mime *mimetype.MIME) bool {
	switch mime.String() {
	case "image/jpeg", "image/gif", "image/png":
		return true
	}
	return false
}

type byteReadCloser struct {
	io.Reader
}

func (brc *byteReadCloser) Close() error {
	brc.Reader = nil
	return nil
}

var ErrTypeNotSupported = errors.New("thumbnail generation not supported for this type")
var ErrEmptyImageGenerated = errors.New("thumbnail generation resulted in empty image (something is broken)")

func Generate(reader io.Reader, mime *mimetype.MIME) (*Thumb, error) {
	if !supportedMime(mime) {
		return nil, ErrTypeNotSupported
	}

	image, _, err := image.Decode(reader)
	if err != nil {
		log.Warnf("error decoding image w/ mime %v: %v", mime.String(), err)
		return nil, err
	}

	resized := resize.Thumbnail(250, 250, image, resize.Lanczos3)

	buf := &bytes.Buffer{}

	switch mime.String() {
	case "image/png", "image/gif":
		mime = mimetype.Lookup("image/png")
		err = png.Encode(buf, resized)
	case "image/jpeg":
		err = jpeg.Encode(buf, resized, nil)
	}

	if err != nil {
		log.Warnf("error encoding image to mime %v: %v", mime.String(), err)
		return nil, err
	}

	bufBytes := buf.Bytes()
	if len(bufBytes) == 0 {
		return nil, ErrEmptyImageGenerated
	}

	return &Thumb{
		ReadCloser: &byteReadCloser{bytes.NewReader(bufBytes)},
		Bytes:      len(bufBytes),
		MIME:       mime,
	}, nil
}
