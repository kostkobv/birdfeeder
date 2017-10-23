package utils

import (
	"bytes"
	"config"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"text/template"
	"unicode/utf16"
)

// Datacoding is a semantic interface for Datacoding enums
type Datacoding string

const (
	// Plain is plain encoding (GSM 7-bit)
	Plain Datacoding = "plain"

	// Unicode is unicode encoding for SMS (UC-2 aka UTF-16)
	Unicode = "unicode"
)

// ErrUC2 is a semantic name for error that is thrown after attempt to encode UC-2 message as GSM 7-bit
var ErrUC2 = errors.New("UC-2")

// UDHEncoder encodes text to string representation of hex values (depends on the used character)
type UDHEncoder interface {
	Encode(m string) *Encoded
	GenerateUDH(p uint8, parts uint8, mesHash uint32) string
	SplitTextMessage(m string) *Encoded
}

type udhenc struct {
	udhUniqueID uint8
	udhCache    map[uint32]uint8
	udhTemplate *template.Template
}

// Encoded is a representation of message split depending on amount of symbols and encoding
type Encoded struct {
	Encoding Datacoding
	Messages []string
}

const udhTemplate = "050003{{.UniqueID}}{{.Parts}}{{.Part}}"

// InitEncoder is UDHEncoder factory method
func InitEncoder() UDHEncoder {
	t, _ := template.New("udh").Parse(udhTemplate) // #nosec

	return &udhenc{
		0x00,
		map[uint32]uint8{},
		t,
	}
}

func getGSM7BitTwoCharsEncodedSymbol(r rune) ([]byte, error) {
	// max capacity - 2 symbols
	result := make([]byte, 0, 2)

	s, ok := config.TwoCharGSMSymbols[r]

	if !ok {
		// it's not 2 space char as well so it's UC-2
		return result, ErrUC2
	}

	result = append(result, config.GSMEscapeSymbol, s)

	return result, nil
}

func getGSM7BitEncodedSymbol(r rune) ([]byte, error) {
	s, ok := config.OneCharGSMSymbols[r]

	if !ok {
		// it's not 1 space char so let's try 2 spaces char
		return getGSM7BitTwoCharsEncodedSymbol(r)
	}

	return []byte{s}, nil
}

func encodeGSM7bit(in string) ([]byte, error) {
	raw := make([]byte, 0, len(in))

	// let's go through the message
	for _, r := range in {

		// encode the symbol
		s, err := getGSM7BitEncodedSymbol(r)

		// it's UC-2! Stopping
		if err != nil {
			return raw, err
		}

		// it's fine, proceed
		raw = append(raw, s...)
	}

	return raw, nil
}

func encodeGSMUC2(in string) []byte {
	r := utf16.Encode([]rune(in))
	buf := new(bytes.Buffer)

	for _, i := range r {
		_ = binary.Write(buf, binary.BigEndian, i) // #nosec
	}

	return buf.Bytes()
}

// by GSM documentation it can be up to 255 but MB documentation says up to 9
const maxSplittedSMSParts = 9

const nonsplittedPlainSMSLength = 160
const splittedPlainSMSLength = 153
const maxPlainSMSLength = splittedPlainSMSLength * maxSplittedSMSParts

const unicodeSymbolLengthBytes = 2
const nonsplittedUnicodeSMSLength = 70
const splittedUnicodeSMSLength = 67
const maxUnicodeSMSLength = splittedUnicodeSMSLength * maxSplittedSMSParts

type smsSplittingLimits struct {
	NonsplittedSMSLength int
	SplittedSMSLength    int
	MaxSMSLength         int
	MaxSMSCharAmount     int
}

func getSMSSplittingLimits(e Datacoding) *smsSplittingLimits {
	var l *smsSplittingLimits

	if e == Plain {
		// message has plain encoding
		l = &smsSplittingLimits{
			NonsplittedSMSLength: nonsplittedPlainSMSLength,
			SplittedSMSLength:    splittedPlainSMSLength,
			MaxSMSLength:         maxPlainSMSLength,
			MaxSMSCharAmount:     maxPlainSMSLength,
		}
	} else {
		// message is unicode
		l = &smsSplittingLimits{
			NonsplittedSMSLength: nonsplittedUnicodeSMSLength * unicodeSymbolLengthBytes,
			SplittedSMSLength:    splittedUnicodeSMSLength * unicodeSymbolLengthBytes,
			MaxSMSLength:         maxUnicodeSMSLength * unicodeSymbolLengthBytes,
			MaxSMSCharAmount:     maxUnicodeSMSLength,
		}
	}

	return l
}

func splitBinaryMessages(enc []byte, e Datacoding) []string {
	s := getSMSSplittingLimits(e)
	l := len(enc)

	// nothing to split here
	if l < s.NonsplittedSMSLength {
		return []string{hex.EncodeToString(enc)}
	}

	var result []string

	// SMS length is longer than it could be. Let's take the max amount of symbols
	if l > s.MaxSMSLength {
		l = s.MaxSMSLength
	}

	// amount of parts
	partsF := float64(l) / float64(s.SplittedSMSLength)
	parts := int(math.Ceil(partsF))

	// iterate through parts
	for i := 0; i < parts; i++ {
		// top range
		up := (i + 1) * s.SplittedSMSLength

		// bottom range
		down := i * s.SplittedSMSLength

		// if it's a last part let's set the top range to the message length, so no panic would be thrown
		if up > l {
			up = l
		}

		// append stringified part to the result
		result = append(result, hex.EncodeToString(enc[down:up]))
	}

	return result
}

// Encode returns the result of hex string encoding of the provided string depending on the used symbols
func (e *udhenc) Encode(m string) *Encoded {
	result := &Encoded{
		Encoding: Plain,
	}

	enc, err := encodeGSM7bit(m)

	if err == ErrUC2 {
		enc = encodeGSMUC2(m)
		result.Encoding = Unicode
	}

	result.Messages = splitBinaryMessages(enc, result.Encoding)

	return result
}

func splitPlainGSM7bit(m string) ([]string, error) {
	result := []string{}

	sum := 0
	part := []rune{}

	for _, r := range m {
		s, err := getGSM7BitEncodedSymbol(r)

		// Unicode!
		if err != nil {
			return result, ErrUC2
		}

		ls := len(s)

		if ls+sum > splittedPlainSMSLength {
			result = append(result, string(part))
			part = []rune{}
			sum = 0
		}

		sum += ls
		part = append(part, r)
	}

	result = append(result, string(part))

	if len(result) == 2 && sum+splittedPlainSMSLength <= nonsplittedPlainSMSLength {
		return []string{m}, nil
	}

	return result, nil
}

func splitPlainGSMUC2(m string) []string {
	runesM := []rune(m)
	l := len(runesM)

	if l <= nonsplittedUnicodeSMSLength {
		return []string{m}
	}

	parts := l/splittedUnicodeSMSLength + 1
	result := []string{}

	for i := 0; i < parts; i++ {
		// top range
		up := (i + 1) * splittedUnicodeSMSLength

		down := i * splittedUnicodeSMSLength

		// if it's a last part let's set the top range to the message length, so no panic would be thrown
		if up > l {
			up = l
		}

		result = append(result, string(runesM[down:up]))
	}

	return result
}

func (e *udhenc) SplitTextMessage(m string) *Encoded {
	result := &Encoded{
		Encoding: Plain,
	}

	var err error

	result.Messages, err = splitPlainGSM7bit(m)

	if err == ErrUC2 {
		result.Messages = splitPlainGSMUC2(m)
		result.Encoding = Unicode
	}

	if len(result.Messages) > maxSplittedSMSParts {
		result.Messages = result.Messages[:maxSplittedSMSParts]
	}

	return result
}

func (e *udhenc) GenerateUDH(p uint8, parts uint8, mesHash uint32) string {
	var uniqueID uint8

	if v, ok := e.udhCache[mesHash]; ok {
		uniqueID = v
	} else {
		e.udhUniqueID++
		uniqueID = e.udhUniqueID
	}

	data := map[string]string{
		"Parts":    e.formatUintString(parts),
		"Part":     e.formatUintString(p),
		"UniqueID": e.formatUintString(uniqueID),
	}

	var udh string

	buf := bytes.NewBufferString(udh)

	_ = e.udhTemplate.Execute(buf, data) // #nosec

	if e.udhUniqueID == 0 {
		e.udhCache = map[uint32]uint8{mesHash: uniqueID}
	} else {
		e.udhCache[mesHash] = uniqueID
	}

	return buf.String()
}

func (e *udhenc) formatUintString(x uint8) string {
	return fmt.Sprintf("%02x", x)
}
