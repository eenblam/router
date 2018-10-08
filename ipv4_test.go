package router

import (
	"testing"
)

func TestMaskFromPrefix(t *testing.T) {
	cases := []struct {
		Name          string
		Input         uint8
		Expected      IPv4
		ExpectedError bool
	}{
		{
			"works mod 8",
			16,
			IPv4{255, 255, 0, 0},
			false,
		},
		{
			"works with borrow bits",
			19,
			IPv4{255, 255, 224, 0},
			false,
		},
		{
			"works with 0",
			0,
			IPv4{0, 0, 0, 0},
			false,
		},
		{
			"works with 32",
			32,
			IPv4{255, 255, 255, 255},
			false,
		},
		{
			"errors above 32",
			33,
			IPv4{0, 0, 0, 0},
			true,
		},
	}
	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			got, err := MaskFromPrefix(test.Input)
			if err != nil {
				if !test.ExpectedError {
					t.Errorf("Unexpected error: %s", err)
				}
				return
			}
			if test.ExpectedError {
				t.Error("Expected error but got none")
				return
			}
			if test.Expected != *got {
				t.Errorf("got %v expected %s", got, test.Expected)
			}
		})
	}
}

func TestIsMask(t *testing.T) {
	cases := []struct {
		Name     string
		Address  IPv4
		Expected bool
	}{
		{
			"/0",
			IPv4{0, 0, 0, 0},
			true,
		},
		{
			"/8",
			IPv4{255, 0, 0, 0},
			true,
		},
		{
			"/16",
			IPv4{255, 255, 0, 0},
			true,
		},
		{
			"/24",
			IPv4{255, 255, 255, 0},
			true,
		},
		{
			"/32",
			IPv4{255, 255, 255, 255},
			true,
		},
		{
			"/15",
			IPv4{255, 254, 0, 0},
			true,
		},
		{
			"/18",
			IPv4{255, 255, 192, 0},
			true,
		},
		{
			"false on leading 0, trailing 1",
			IPv4{64, 0, 0, 0},
			false,
		},
		{
			"false on leading 0, last bit 1",
			IPv4{255, 1, 0, 0},
			false,
		},
		{
			"false on 101",
			IPv4{255, 160, 0, 0},
			false,
		},
	}
	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			got := test.Address.IsMask()
			if test.Expected != got {
				t.Errorf("expected %t got %t for %v", test.Expected, got, test.Address)
			}
		})
	}
}
func TestMaskWith(t *testing.T) {
	cases := []struct {
		Name     string
		Address  IPv4
		Mask     IPv4
		Expected *IPv4
	}{
		{
			"Simple test",
			IPv4{192, 168, 0, 0},
			IPv4{255, 255, 0, 0},
			&IPv4{192, 168, 0, 0},
		},
		{
			"works with a /16",
			IPv4{192, 168, 0, 0},
			IPv4{255, 255, 0, 0},
			&IPv4{192, 168, 0, 0},
		},
		{
			"works with a /18",
			IPv4{192, 168, 64, 0},
			IPv4{255, 255, 192, 0},
			&IPv4{192, 168, 64, 0},
		},
		{
			"nil on /18 mismatch",
			IPv4{192, 168, 0, 0},
			IPv4{255, 255, 64, 0},
			nil,
		},
	}
	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			masked := test.Address.MaskWith(test.Mask)
			if test.Expected == nil && masked == nil {
				return
			} else if test.Expected == nil && masked != nil {
				t.Errorf("Expected nil but got %v", *masked)
			} else if masked == nil {
				t.Errorf("Expected %v but got nil", *test.Expected)
			} else if *test.Expected != *masked {
				t.Errorf("got %v expected %v", *masked, *test.Expected)
			}
		})
	}
}
