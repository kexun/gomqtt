package packet

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentifiedDecode(t *testing.T) {
	pktBytes := []byte{
		byte(PUBACK << 4),
		2,
		0, // packet ID MSB
		7, // packet ID LSB
	}

	n, pid, err := identifiedDecode(pktBytes, PUBACK)
	assert.NoError(t, err)
	assert.Equal(t, 4, n)
	assert.Equal(t, ID(7), pid)
}

func TestIdentifiedDecodeError1(t *testing.T) {
	pktBytes := []byte{
		byte(PUBACK << 4),
		1, // < wrong remaining length
		0, // packet ID MSB
		7, // packet ID LSB
	}

	n, pid, err := identifiedDecode(pktBytes, PUBACK)
	assert.Error(t, err)
	assert.Equal(t, 2, n)
	assert.Equal(t, ID(0), pid)
}

func TestIdentifiedDecodeError2(t *testing.T) {
	pktBytes := []byte{
		byte(PUBACK << 4),
		2,
		7, // packet ID LSB
		// < insufficient bytes
	}

	n, pid, err := identifiedDecode(pktBytes, PUBACK)
	assert.Error(t, err)
	assert.Equal(t, 2, n)
	assert.Equal(t, ID(0), pid)
}

func TestIdentifiedDecodeError3(t *testing.T) {
	pktBytes := []byte{
		byte(PUBACK << 4),
		2,
		0, // packet ID LSB
		0, // packet ID MSB < zero id
	}

	n, pid, err := identifiedDecode(pktBytes, PUBACK)
	assert.Error(t, err)
	assert.Equal(t, 4, n)
	assert.Equal(t, ID(0), pid)
}

func TestIdentifiedEncode(t *testing.T) {
	pktBytes := []byte{
		byte(PUBACK << 4),
		2,
		0, // packet ID MSB
		7, // packet ID LSB
	}

	dst := make([]byte, identifiedLen())
	n, err := identifiedEncode(dst, 7, PUBACK)

	assert.NoError(t, err)
	assert.Equal(t, 4, n)
	assert.Equal(t, pktBytes, dst[:n])
}

func TestIdentifiedEncodeError1(t *testing.T) {
	dst := make([]byte, 3) // < insufficient buffer
	n, err := identifiedEncode(dst, 7, PUBACK)

	assert.Error(t, err)
	assert.Equal(t, 0, n)
}

func TestIdentifiedEncodeError2(t *testing.T) {
	dst := make([]byte, identifiedLen())
	n, err := identifiedEncode(dst, 0, PUBACK) // < zero id

	assert.Error(t, err)
	assert.Equal(t, 0, n)
}

func TestIdentifiedEqualDecodeEncode(t *testing.T) {
	pktBytes := []byte{
		byte(PUBACK << 4),
		2,
		0, // packet ID MSB
		7, // packet ID LSB
	}

	pkt := &Puback{}
	n, err := pkt.Decode(pktBytes)

	assert.NoError(t, err)
	assert.Equal(t, 4, n)

	dst := make([]byte, 100)
	n2, err := identifiedEncode(dst, 7, PUBACK)

	assert.NoError(t, err)
	assert.Equal(t, 4, n2)
	assert.Equal(t, pktBytes, dst[:n2])

	n3, pid, err := identifiedDecode(pktBytes, PUBACK)
	assert.NoError(t, err)
	assert.Equal(t, 4, n3)
	assert.Equal(t, ID(7), pid)
}

func BenchmarkIdentifiedEncode(b *testing.B) {
	pkt := &Puback{}
	pkt.ID = 1

	buf := make([]byte, pkt.Len())

	for i := 0; i < b.N; i++ {
		_, err := pkt.Encode(buf)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkIdentifiedDecode(b *testing.B) {
	pktBytes := []byte{
		byte(PUBACK << 4),
		2,
		0, // packet ID MSB
		1, // packet ID LSB
	}

	pkt := &Puback{}

	for i := 0; i < b.N; i++ {
		_, err := pkt.Decode(pktBytes)
		if err != nil {
			panic(err)
		}
	}
}

func testIdentifiedImplementation(t *testing.T, pkt Generic) {
	assert.Equal(t, fmt.Sprintf("<%s ID=1>", pkt.Type().String()), pkt.String())

	buf := make([]byte, pkt.Len())
	n, err := pkt.Encode(buf)
	assert.NoError(t, err)
	assert.Equal(t, 4, n)

	n, err = pkt.Decode(buf)
	assert.NoError(t, err)
	assert.Equal(t, 4, n)
}

func TestPubackImplementation(t *testing.T) {
	pkt := NewPuback()
	pkt.ID = 1

	testIdentifiedImplementation(t, pkt)
}

func TestPubcompImplementation(t *testing.T) {
	pkt := NewPubcomp()
	pkt.ID = 1

	testIdentifiedImplementation(t, pkt)
}

func TestPubrecImplementation(t *testing.T) {
	pkt := NewPubrec()
	pkt.ID = 1

	testIdentifiedImplementation(t, pkt)
}

func TestPubrelImplementation(t *testing.T) {
	pkt := NewPubrel()
	pkt.ID = 1

	testIdentifiedImplementation(t, pkt)
}

func TestUnsubackImplementation(t *testing.T) {
	pkt := NewUnsuback()
	pkt.ID = 1

	testIdentifiedImplementation(t, pkt)
}
