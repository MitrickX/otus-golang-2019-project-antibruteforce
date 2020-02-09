package entities

import (
	"testing"
)

func TestIP_NewSubNetIP(t *testing.T) {
	ip, err := New("127.0.0.1")
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	expectedIP := IP("127.0.0.1")
	if ip != expectedIP {
		t.Fatalf("unexpected ip `%s` instreadof `%s`", ip, expectedIP)
	}
}

func TestIP_NewHostIP(t *testing.T) {
	ip, err := New("127.0.0.0/24")
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	expectedIP := IP("127.0.0.0/24")
	if ip != expectedIP {
		t.Fatalf("unexpected ip `%s` instreadof `%s`", ip, expectedIP)
	}
}

func TestIP_NewInvalidIP(t *testing.T) {
	_, err := New("adfasdf")
	if err == nil {
		t.Fatal("expect error")
	}
}

func TestIP_NewWithMaskPartSubNetIP(t *testing.T) {
	ip, err := NewWithMaskPart("127.0.0.0/24")
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	expectedIP := IP("127.0.0.0/24")
	if ip != expectedIP {
		t.Fatalf("unexpected ip `%s` instreadof `%s`", ip, expectedIP)
	}
}

func TestIP_NewWithMaskPartHostIP(t *testing.T) {
	_, err := NewWithMaskPart("127.0.0.1")
	if err == nil {
		t.Fatal("expect error")
	}
}

func TestIP_NewWithMaskPartInvalidIP(t *testing.T) {
	_, err := NewWithMaskPart("dfasdf")
	if err == nil {
		t.Fatal("expect error")
	}
}

func TestIP_NewWithoutMaskPartHostIP(t *testing.T) {
	ip, err := NewWithoutMaskPart("127.0.0.1")
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	expectedIP := IP("127.0.0.1")
	if ip != expectedIP {
		t.Fatalf("unexpected ip `%s` instreadof `%s`", ip, expectedIP)
	}
}

func TestIP_NewWithoutMaskPartSubNetIP(t *testing.T) {
	_, err := NewWithoutMaskPart("127.0.0.0/24")
	if err == nil {
		t.Fatal("expect error")
	}
}

func TestIP_NewWithoutMaskPartInvalidIP(t *testing.T) {
	_, err := NewWithoutMaskPart("dfasdf")
	if err == nil {
		t.Fatal("expect error")
	}
}

func TestIP_DropMaskPartSubNetIP(t *testing.T) {
	ip := IP("127.0.0.0/24")
	expectedIP := IP("127.0.0.0")
	resIP := ip.DropMaskPart()

	if expectedIP != resIP {
		t.Fatalf("unexpected ip `%s` instreadof `%s`", resIP, expectedIP)
	}
}

func TestIP_DropMaskPartHostIP(t *testing.T) {
	ip := IP("127.0.0.1")
	expectedIP := IP("127.0.0.1")
	resIP := ip.DropMaskPart()

	if expectedIP != resIP {
		t.Fatalf("unexpected ip `%s` instreadof `%s`", resIP, expectedIP)
	}
}

func TestIP_HasMaskPartSubNetIP(t *testing.T) {
	ip := IP("127.0.0.0/24")
	if !ip.HasMaskPart() {
		t.Fatalf("unexpected false")
	}
}

func TestIP_HasMaskPartHostIP(t *testing.T) {
	ip := IP("127.0.0.1")
	if ip.HasMaskPart() {
		t.Fatalf("unexpected true")
	}
}

func TestIP_Parse(t *testing.T) {
	ip := IP("127.0.0.1")

	netIP := ip.Parse()
	if netIP.String() != string(ip) {
		t.Fatalf("unexpected result of conveerting net.IPBits `%#v` to string `%s` instreadof `%s`",
			netIP,
			netIP.String(),
			ip,
		)
	}
}

func TestIP_ParseAsCIDRSubNetIP(t *testing.T) {
	ip := IP("127.0.0.0/24")

	netIP, maskIP, err := ip.ParseAsCIDR()
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	if netIP.String() != string(ip.DropMaskPart()) {
		t.Fatalf("unexpected that netIP %s not equals ip after drop mask %s",
			netIP.String(),
			ip.DropMaskPart(),
		)
	}

	if maskIP.IP.String() != string(ip.DropMaskPart()) {
		t.Fatalf("unexpected that maskIP.IPBits %s not equals ip after drop mask %s",
			netIP.String(),
			ip.DropMaskPart(),
		)
	}

	ones, _ := maskIP.Mask.Size()

	//nolint:gomnd
	if ones != 24 {
		t.Fatalf("unexpected len %d instreadof %d", ones, 24)
	}
}

func TestIP_ParseAsCIDRHostIP(t *testing.T) {
	ip := IP("127.0.0.0")

	_, _, err := ip.ParseAsCIDR()
	if err == nil {
		t.Fatal("unexpected that there is no error")
	}
}

func TestIP_IsConform(t *testing.T) {
	subNetIP := IP("127.0.0.0/24")

	expected := map[IP]bool{
		IP("127.0.0.1"):   true,
		IP("127.0.0.50"):  true,
		IP("127.0.0.255"): true,
		IP("128.0.0.3"):   false,
	}

	for ip, expectedRes := range expected {
		res := subNetIP.IsConform(ip)
		if expectedRes != res {
			t.Fatalf("unexpected conform result %t for `%s`", res, ip)
		}
	}
}
