package heic

import (
	"bufio"
	"bytes"
	"fmt"
	"image/jpeg"
	"io"

	"github.com/dsoprea/go-exif/v3"
	"github.com/dsoprea/go-heic-exif-extractor/v2"
	"github.com/dsoprea/go-jpeg-image-structure/v2"
	"github.com/strukturag/libheif-go"
)

// Read HEIC image and EXIF data from 'r' and write it back as JPEG and EXIF data to 'wr'.
func ToJPEG(r io.Reader, wr io.Writer) error {

	heic_body, err := io.ReadAll(r)

	if err != nil {
		return fmt.Errorf("Failed to read input body, %w", err)
	}

	// First decode the HEIC image

	im_ctx, err := libheif.NewContext()

	if err != nil {
		return fmt.Errorf("Failed to create new libheif context, %w", err)
	}

	err = im_ctx.ReadFromMemory(heic_body)

	if err != nil {
		return fmt.Errorf("Failed to read input data, %w", err)
	}

	im_handle, err := im_ctx.GetPrimaryImageHandle()

	if err != nil {
		return fmt.Errorf("Failed to derive primary image handler, %w", err)
	}

	h_im, err := im_handle.DecodeImage(libheif.ColorspaceUndefined, libheif.ChromaUndefined, nil)

	if err != nil {
		return fmt.Errorf("Failed to decode image, %w", err)
	}

	// Convert to image.Image and then write as JPEG file

	im, err := h_im.GetImage()

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	var jpeg_buf bytes.Buffer
	jpeg_wr := bufio.NewWriter(&jpeg_buf)

	jpeg_opts := &jpeg.Options{
		Quality: 100,
	}

	err = jpeg.Encode(jpeg_wr, im, jpeg_opts)

	if err != nil {
		return fmt.Errorf("Failed to write JPEG data, %w", err)
	}

	jpeg_wr.Flush()

	// Parse EXIF data out of HEIC image

	hemp := new(heicexif.HeicExifMediaParser)

	mc, err := hemp.ParseBytes(heic_body)

	if err != nil {
		return fmt.Errorf("Failed to parse EXIF data from input, %w", err)
	}

	rootIfd, _, err := mc.Exif()

	if err != nil {
		return fmt.Errorf("Failed to derive EXIF data from input, %w", err)
	}

	rootIb := exif.NewIfdBuilderFromExistingChain(rootIfd)

	// Finally write JPEG data and EXIF back to 'wr'

	jmp := jpegstructure.NewJpegMediaParser()

	intfc, err := jmp.ParseBytes(jpeg_buf.Bytes())

	if err != nil {
		return fmt.Errorf("Failed to parse JPEG data, %w", err)
	}

	sl := intfc.(*jpegstructure.SegmentList)

	err = sl.SetExif(rootIb)

	if err != nil {
		return fmt.Errorf("Failed to assign EXIF to JPEG data, %w", err)
	}

	err = sl.Write(wr)

	if err != nil {
		return fmt.Errorf("Failed to write JPEG+EXIF data, %w", err)
	}

	return nil
}
