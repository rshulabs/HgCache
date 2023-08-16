package discovery

import "testing"

func TestUserRegistry1(t *testing.T) {
	u := NewUserService("127.0.0.1:4543")
	UserRegisty(u)
}

func TestUserRegistry2(t *testing.T) {
	u := NewUserService("127.0.0.1:4544")
	UserRegisty(u)
}

func TestUserDiscovery1(t *testing.T) {
	UserDiscovery("user")
}

func TestUserDiscovery2(t *testing.T) {
	UserDiscovery("user")
}

func TestUserDiscovery3(t *testing.T) {
	UserDiscovery("user")
}