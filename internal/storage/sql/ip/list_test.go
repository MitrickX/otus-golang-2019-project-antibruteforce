package ip

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/domain/entities"
)

func (row Row) String() string {
	maskBitsStr := "<nil>"
	if row.MaskBits != nil {
		maskBitsStr = string(*row.MaskBits)
	}

	maskIntStr := "<nil>"
	if row.MaskInt != nil {
		maskIntStr = fmt.Sprintf("%d", *row.MaskInt)
	}

	return fmt.Sprintf("{\n  IPBits: %s\n  MaskBits: %s\n  IP: %s\n  MaskInt: %s\n  IsV6: %t\n  Kind: %s\n}",
		row.IPBits, maskBitsStr, row.IP, maskIntStr, row.IsV6, row.Kind)
}

func TestConvertByteSliceToBitMask(t *testing.T) {
	bs := []byte{127, 0, 0, 1}
	bitMask := convertByteSliceToBitMask(bs)
	expectedBitMask := BitMask("01111111000000000000000000000001")

	if bitMask != expectedBitMask {
		t.Fatalf("unexpected bit mask `%s` instead of `%s` when convert byte slice %s",
			bitMask, expectedBitMask, bs)
	}
}

func TestConvertIpToRow_IPv4(t *testing.T) {
	ip := entities.IP("127.0.0.1")

	row, err := convertIPToRow(ip, "black")
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	expectedRow := Row{
		IPBits:   BitMask("01111111000000000000000000000001"),
		MaskBits: nil,
		IP:       ip,
		MaskInt:  nil,
		IsV6:     false,
		Kind:     "black",
	}

	assertRowsEquals(t, expectedRow, row)
}

func TestConvertIpToRow_SubNetIPv4(t *testing.T) {
	ip := entities.IP("127.0.0.1/24")

	row, err := convertIPToRow(ip, "black")
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	maskBits := BitMask("11111111111111111111111100000000")
	maskInt := 24

	expectedRow := Row{
		IPBits:   BitMask("01111111000000000000000000000001"),
		MaskBits: &maskBits,
		IP:       ip,
		MaskInt:  &maskInt,
		IsV6:     false,
		Kind:     "black",
	}

	assertRowsEquals(t, expectedRow, row)
}

func TestConvertIpToRow_IPv6(t *testing.T) {
	ip := entities.IP("2001:DB8:0:1234::2")

	row, err := convertIPToRow(ip, "black")
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	bitMask := BitMask(
		"00100000" +
			"00000001" +
			"00001101" +
			"10111000" +
			"00000000" +
			"00000000" +
			"00010010" +
			"00110100" +
			"00000000" +
			"00000000" +
			"00000000" +
			"00000000" +
			"00000000" +
			"00000000" +
			"00000000" +
			"00000010")

	expectedRow := Row{
		IPBits:   bitMask,
		MaskBits: nil,
		IP:       ip,
		MaskInt:  nil,
		IsV6:     true,
		Kind:     "black",
	}

	assertRowsEquals(t, expectedRow, row)
}

func TestConvertIpToRow_SubNetIPv6(t *testing.T) {
	ip := entities.IP("2001:DB8:0:1234::2/64")

	row, err := convertIPToRow(ip, "black")
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	maskBits := BitMask(
		"11111111" +
			"11111111" +
			"11111111" +
			"11111111" +
			"11111111" +
			"11111111" +
			"11111111" +
			"11111111" +
			"00000000" +
			"00000000" +
			"00000000" +
			"00000000" +
			"00000000" +
			"00000000" +
			"00000000" +
			"00000000")

	maskInt := 64

	ipBits := BitMask(
		"00100000" +
			"00000001" +
			"00001101" +
			"10111000" +
			"00000000" +
			"00000000" +
			"00010010" +
			"00110100" +
			"00000000" +
			"00000000" +
			"00000000" +
			"00000000" +
			"00000000" +
			"00000000" +
			"00000000" +
			"00000010")

	expectedRow := Row{
		IPBits:   ipBits,
		MaskBits: &maskBits,
		IP:       ip,
		MaskInt:  &maskInt,
		IsV6:     true,
		Kind:     "black",
	}

	assertRowsEquals(t, expectedRow, row)
}

func assertRowsEquals(t *testing.T, expectedRow Row, row Row) {
	if !reflect.DeepEqual(expectedRow, row) {
		t.Fatalf("unexpected row:\n%s\ninstead of:\n%s", row, expectedRow)
	}
}
