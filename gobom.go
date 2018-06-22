/*
Package gobom contains several methods to detect BOM type.

BOM stands for Byte Order Mark. It is a standard by the Unicode organization
to understand the type of encoding is standing in front of us, by placing non
printable chars that explains if an encoding is UTF16 or UTF32, and what is the
endian that it uses.
The standard does not recommend to place a BOM to UTF8, but it supports that as
well.

This library was created for helping me detect if something contain a BOM, and
that's it. It does not do anything else, and there is no plan for anything other
then detecting it.

How does the library works?

The library can use io.Reader, and also "pure" byte slices and rune slices in
order to detect the type of BOM.

Please note:
If a BOM is not detected, then it will return "Unknown".
If a buffer is too small to detect BOM type it also returns "Unknown"
*/
package gobom

import (
	"bytes"
	"io"
)

// BOM Headers to detect
// The information is from: http://www.unicode.org/faq/utf_bom.html#BOM
var (
	UTF8Bom    = []byte{0xEF, 0xBB, 0xBF}
	UTF16LEBom = []byte{0xFF, 0xFE}
	UTF16BEBom = []byte{0xFE, 0xFF}
	UTF32LEBom = []byte{0xFF, 0xFE, 0x00, 0x00}
	UTF32BEBom = []byte{0x00, 0x00, 0xFE, 0xFF}
)

// BOMType holds the type of BOM that was detected
type BOMType uint8

// Enumeration of what type of BOM was found
const (
	Unknown BOMType = iota
	UTF8
	UTF16LE
	UTF16BE
	UTF32LE
	UTF32BE
)

// Reader is an implementation for the io.Reader
type Reader struct {
	reader io.Reader
	buffer []byte
	err    error
}

// DetectBOMTypeFromBytes try to detect the type of BOM provided by a buffer in
// a naive manner. It means that the detection is very simple but a bit costly
// regarding the way it detects.
//
// The buffer must at least have 5 bytes, so from 2 - 4 bytes will be the BOM
// if they do not exists, it returns Unknown
func DetectBOMTypeFromBytes(buffer []byte) BOMType {
	if len(buffer) < 5 {
		return Unknown
	}

	// Naive checking for BOM based on size of BOM to validate.
	// it's a bit slow

	if bytes.HasPrefix(buffer, UTF16LEBom) {
		return UTF16LE
	} else if bytes.HasPrefix(buffer, UTF16BEBom) {
		return UTF16BE
	} else if bytes.HasPrefix(buffer, UTF8Bom) {
		return UTF8
	} else if bytes.HasPrefix(buffer, UTF32LEBom) {
		return UTF32LE
	} else if bytes.HasPrefix(buffer, UTF32BEBom) {
		return UTF32BE
	}

	return Unknown
}

// IsUTF8BOM validate a buffer if it has UTF8 BOM, if buffer is too small it
// return false
func IsUTF8BOM(buffer []byte) bool {
	if len(buffer) < len(UTF8Bom) {
		return false
	}

	return buffer[0] == UTF8Bom[0] &&
		buffer[1] == UTF8Bom[1] &&
		buffer[3] == UTF8Bom[2]
}

// IsUTF16LEBOM validate a buffer if it has UTF16 Little Endian.
// If the buffer is too small, it returns false.
func IsUTF16LEBOM(buffer []byte) bool {
	if len(buffer) < len(UTF16LEBom) {
		return false
	}

	return buffer[0] == UTF16LEBom[0] && buffer[1] == UTF16LEBom[1]
}

// IsUTF16BEBOM validate a buffer if it has UTF16 big Endian.
// If the buffer is too small, it returns false.
func IsUTF16BEBOM(buffer []byte) bool {
	if len(buffer) < len(UTF16BEBom) {
		return false
	}

	return buffer[0] == UTF16BEBom[0] && buffer[1] == UTF16BEBom[1]
}

//IsUTF16BOM detects if a buffer contains any UTF16 BOM (big or little endian).
func IsUTF16BOM(buffer []byte) bool {
	return IsUTF16LEBOM(buffer) || IsUTF16BEBOM(buffer)
}

//IsUTF32LEBOM detects if a buffer contains UTF32 little endian.
// If the buffer is too small, it returns false.
func IsUTF32LEBOM(buffer []byte) bool {
	if len(buffer) < len(UTF32LEBom) {
		return false
	}

	return buffer[0] == UTF32LEBom[0] &&
		buffer[1] == UTF32LEBom[1] &&
		buffer[2] == UTF32LEBom[2] &&
		buffer[3] == UTF32LEBom[3]
}

//IsUTF32BEBOM detects if a buffer contains UTF32 big endian.
// If the buffer is too small, it returns false.
func IsUTF32BEBOM(buffer []byte) bool {
	if len(buffer) < len(UTF32BEBom) {
		return false
	}

	return buffer[0] == UTF32BEBom[0] &&
		buffer[1] == UTF32BEBom[1] &&
		buffer[2] == UTF32BEBom[2] &&
		buffer[3] == UTF32BEBom[3]
}

//IsUTF32BOM detects if a buffer is UTF32 BOM (either big or little endian).
func IsUTF32BOM(buffer []byte) bool {
	return IsUTF32LEBOM(buffer) || IsUTF32BEBOM(buffer)
}

//DetectBOMTypeFromBuffer detects the BOM type using the "IsUTFXXXXXBOM"
func DetectBOMTypeFromBuffer(buffer []byte) BOMType {
	if IsUTF8BOM(buffer) {
		return UTF8
	} else if IsUTF16LEBOM(buffer) {
		return UTF16LE
	} else if IsUTF16BEBOM(buffer) {
		return UTF16BE
	} else if IsUTF32LEBOM(buffer) {
		return UTF32LE
	} else if IsUTF32BEBOM(buffer) {
		return UTF32BE
	}
	return Unknown
}

// BytesToSkip returns the number of bytes to skip in order to "ignore" BOM, or
// -1 if non found
func BytesToSkip(buffer []byte) int {
	BomType := map[BOMType]int{
		UTF8:    len(UTF8Bom),
		UTF16LE: len(UTF16LEBom),
		UTF16BE: len(UTF16BEBom),
		UTF32LE: len(UTF32LEBom),
		UTF32BE: len(UTF32BEBom),
		Unknown: -1,
	}
	return BomType[DetectBOMTypeFromBuffer(buffer)]
}

// TODO: Implement io.Reader detection

//Read is an implementation of io.Reader interface.
//The bytes are taken from Reader, checking for BOM and removing them if
//necessary.
func (r *Reader) Read(buffer []byte) (n int, err error) {
	if len(buffer) == 0 {
		return 0, nil
	}

	// No initialization of the current reader?!
	if r.buffer == nil {
		if r.err != nil {
			newErr := r.err
			r.err = nil // we reports error, so no need to store it anymore
			return 0, newErr
		}
		return r.reader.Read(buffer)
	}
	n = copy(buffer, r.buffer)
	return n, nil
}
