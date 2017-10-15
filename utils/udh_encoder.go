package utils

import (
	"bytes"
	"config"
	"encoding/binary"
	"encoding/hex"
	"errors"
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
var ErrUC2 error = errors.New("UC-2")

// UDHEncoder encodes text to string representation of hex values (depends on the used character)
type UDHEncoder interface {
	Encode(m string) (*Encoded, error)
}

type udhenc struct {
}

// Encoded is a representation of message split depending on amount of symbols and encoding
type Encoded struct {
	Encoding Datacoding
	Messages []string
}

// InitEncoder is UDHEncoder factory method
func InitEncoder() UDHEncoder {
	return &udhenc{}
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

func splitMessages(enc []byte) []string {
	return []string{hex.EncodeToString(enc)}
}

func octetsAmount(n, block int) int {
	b := n / block

	if n%block != 0 {
		b++
	}

	return b
}

const blockLength = 8
const charLength = 7

func packOctet(out []byte, b byte, octIndex int, bitIndex uint8) (int, uint8) {
	var i uint8

	// let's go through the char slice
	for ; i < charLength; i++ {
		// convert 8 bits to 7-bit form
		out[octIndex] = out[octIndex] | b>>i&1<<bitIndex

		// moving to the next bit in the block
		bitIndex++

		if bitIndex == blockLength {
			// moving to the next octet
			octIndex++
			// reset the bit index
			bitIndex = 0
		}
	}

	return octIndex, bitIndex
}

func packOctets(raw []byte) []byte {
	octets := make([]byte, octetsAmount(len(raw)*charLength, blockLength))

	var octIndex int   // current octet in octets
	var bitIndex uint8 // current bit index in octet
	var b byte         // current byte in octet
	for i := range raw {
		b = raw[i]
		octIndex, bitIndex = packOctet(octets, b, octIndex, bitIndex)
	}

	return handleOctetsEnding(octets, bitIndex, octIndex, b)
}

func handleOctetsEnding(octets []byte, bitIndex uint8, octIndex int, b byte) []byte {
	// 7 zero-bits could be confused with @ so <CR> is added in that case
	if blockLength-bitIndex == charLength {
		packOctet(octets, config.GSMCRSymbol, octIndex, bitIndex)
	} else if bitIndex == 0 && b == config.GSMCRSymbol {
		// if data ends with <CR> we will append empty octet and then add another <CR>
		octets = append(octets, config.EmptyOctet)
		packOctet(octets, config.GSMCRSymbol, octIndex, bitIndex)
	}

	return octets
}

// Encode returns the result of hex string encoding of the provided string depending on the used symbols
func (e *udhenc) Encode(m string) (*Encoded, error) {
	result := &Encoded{
		Encoding: Plain,
	}

	enc, err := encodeGSM7bit(m)

	if err == ErrUC2 {
		enc = encodeGSMUC2(m)
		result.Encoding = Unicode
	} else {
		enc = packOctets(enc)
	}

	result.Messages = splitMessages(enc)

	return result, nil
}
