package entities

import (
	"testing"
)

func TestIP_DropMaskPart1(t *testing.T) {
	ip := IP("127.0.0.0/24")
	expectedIP := IP("127.0.0.0")
	resIP := ip.DropMaskPart()
	if expectedIP != resIP {
		t.Fatalf("unexpected ip `%s` instreadof `%s`", resIP, expectedIP)
	}
}

func TestIP_DropMaskPart2(t *testing.T) {
	ip := IP("127.0.0.1")
	expectedIP := IP("127.0.0.1")
	resIP := ip.DropMaskPart()
	if expectedIP != resIP {
		t.Fatalf("unexpected ip `%s` instreadof `%s`", resIP, expectedIP)
	}
}

func TestIP_HasMaskPart1(t *testing.T) {
	ip := IP("127.0.0.0/24")
	if !ip.HasMaskPart() {
		t.Fatalf("unexpected false")
	}
}

func TestIP_HasMaskPart2(t *testing.T) {
	ip := IP("127.0.0.1")
	if ip.HasMaskPart() {
		t.Fatalf("unexpected true")
	}
}

func TestIP_Parse(t *testing.T) {
	ip := IP("127.0.0.1")
	netIP := ip.Parse()
	if netIP.String() != string(ip) {
		t.Fatalf("unexpected result of conveerting net.IP `%#v` to string `%s` instreadof `%s`",
			netIP,
			netIP.String(),
			ip,
		)
	}
}

func TestIP_ParseAsCIDR1(t *testing.T) {
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
		t.Fatalf("unexpected that maskIP.IP %s not equals ip after drop mask %s",
			netIP.String(),
			ip.DropMaskPart(),
		)
	}

	ones, _ := maskIP.Mask.Size()

	if ones != 24 {
		t.Fatalf("unexpected len %d instreadof %d", ones, 24)
	}
}

func TestIP_ParseAsCIDR2(t *testing.T) {
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
