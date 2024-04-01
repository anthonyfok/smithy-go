package cbor

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestDecode_InvalidArgument(t *testing.T) {
	cases := []errorTest{}
	for name, c := range map[string]struct {
		In  []byte
		Err string
	}{
		"uint/1": {
			[]byte{0<<5 | 24},
			"arg len 1 greater than remaining buf len",
		},
		"uint/2": {
			[]byte{0<<5 | 25, 0},
			"arg len 2 greater than remaining buf len",
		},
		"uint/4": {
			[]byte{0<<5 | 26, 0, 0, 0},
			"arg len 4 greater than remaining buf len",
		},
		"uint/8": {
			[]byte{0<<5 | 27, 0, 0, 0, 0, 0, 0, 0},
			"arg len 8 greater than remaining buf len",
		},
		"uint/?": {
			[]byte{0<<5 | 31},
			"unexpected minor value 31",
		},
		"negint/1": {
			[]byte{1<<5 | 24},
			"arg len 1 greater than remaining buf len",
		},
		"negint/2": {
			[]byte{1<<5 | 25, 0},
			"arg len 2 greater than remaining buf len",
		},
		"negint/4": {
			[]byte{1<<5 | 26, 0, 0, 0},
			"arg len 4 greater than remaining buf len",
		},
		"negint/8": {
			[]byte{1<<5 | 27, 0, 0, 0, 0, 0, 0, 0},
			"arg len 8 greater than remaining buf len",
		},
		"negint/?": {
			[]byte{1<<5 | 31},
			"unexpected minor value 31",
		},
		"slice/1": {
			[]byte{2<<5 | 24},
			"arg len 1 greater than remaining buf len",
		},
		"slice/2": {
			[]byte{2<<5 | 25, 0},
			"arg len 2 greater than remaining buf len",
		},
		"slice/4": {
			[]byte{2<<5 | 26, 0, 0, 0},
			"arg len 4 greater than remaining buf len",
		},
		"slice/8": {
			[]byte{2<<5 | 27, 0, 0, 0, 0, 0, 0, 0},
			"arg len 8 greater than remaining buf len",
		},
		"string/1": {
			[]byte{3<<5 | 24},
			"arg len 1 greater than remaining buf len",
		},
		"string/2": {
			[]byte{3<<5 | 25, 0},
			"arg len 2 greater than remaining buf len",
		},
		"string/4": {
			[]byte{3<<5 | 26, 0, 0, 0},
			"arg len 4 greater than remaining buf len",
		},
		"string/8": {
			[]byte{3<<5 | 27, 0, 0, 0, 0, 0, 0, 0},
			"arg len 8 greater than remaining buf len",
		},
		"list/1": {
			[]byte{4<<5 | 24},
			"arg len 1 greater than remaining buf len",
		},
		"list/2": {
			[]byte{4<<5 | 25, 0},
			"arg len 2 greater than remaining buf len",
		},
		"list/4": {
			[]byte{4<<5 | 26, 0, 0, 0},
			"arg len 4 greater than remaining buf len",
		},
		"list/8": {
			[]byte{4<<5 | 27, 0, 0, 0, 0, 0, 0, 0},
			"arg len 8 greater than remaining buf len",
		},
		"map/1": {
			[]byte{5<<5 | 24},
			"arg len 1 greater than remaining buf len",
		},
		"map/2": {
			[]byte{5<<5 | 25, 0},
			"arg len 2 greater than remaining buf len",
		},
		"map/4": {
			[]byte{5<<5 | 26, 0, 0, 0},
			"arg len 4 greater than remaining buf len",
		},
		"map/8": {
			[]byte{5<<5 | 27, 0, 0, 0, 0, 0, 0, 0},
			"arg len 8 greater than remaining buf len",
		},
		"tag/1": {
			[]byte{6<<5 | 24},
			"arg len 1 greater than remaining buf len",
		},
		"tag/2": {
			[]byte{6<<5 | 25, 0},
			"arg len 2 greater than remaining buf len",
		},
		"tag/4": {
			[]byte{6<<5 | 26, 0, 0, 0},
			"arg len 4 greater than remaining buf len",
		},
		"tag/8": {
			[]byte{6<<5 | 27, 0, 0, 0, 0, 0, 0, 0},
			"arg len 8 greater than remaining buf len",
		},
		"tag/?": {
			[]byte{6<<5 | 31},
			"unexpected minor value 31",
		},
		"major7/float16": {
			[]byte{7<<5 | 25, 0},
			"incomplete float16 at end of buf",
		},
		"major7/float32": {
			[]byte{7<<5 | 26, 0, 0, 0},
			"incomplete float32 at end of buf",
		},
		"major7/float64": {
			[]byte{7<<5 | 27, 0, 0, 0, 0, 0, 0, 0},
			"incomplete float64 at end of buf",
		},
		"major7/?": {
			[]byte{7<<5 | 31},
			"unexpected minor value 31",
		},
	} {
		t.Run(name, func(t *testing.T) {
			cases = append(cases, errorTest{
				Description: fmt.Sprintf("TestDecode_InvalidArgument - %s - %s", name, c.Err),
				Input:       hex.EncodeToString(c.In),
			})

			_, _, err := decode(c.In)
			if err == nil {
				t.Errorf("expect err %s", c.Err)
			}
			if aerr := err.Error(); !strings.Contains(aerr, c.Err) {
				t.Errorf("expect err %s, got %s", c.Err, aerr)
			}
		})
	}

	dumpErrorCases(t, "TestDecodeError_InvalidArgument", cases)
}

func TestDecode_InvalidSlice(t *testing.T) {
	cases := []errorTest{}
	for name, c := range map[string]struct {
		In  []byte
		Err string
	}{
		"slice/1, not enough bytes": {
			[]byte{2<<5 | 24, 1},
			"slice len 1 greater than remaining buf len",
		},
		"slice/?, no break": {
			[]byte{2<<5 | 31},
			"expected break marker",
		},
		"slice/?, invalid nested major": {
			[]byte{2<<5 | 31, 3<<5 | 0},
			"unexpected major type 3 in indefinite slice",
		},
		"slice/?, nested indefinite": {
			[]byte{2<<5 | 31, 2<<5 | 31},
			"nested indefinite slice",
		},
		"slice/?, invalid nested definite": {
			[]byte{2<<5 | 31, 2<<5 | 24, 1},
			"decode subslice: slice len 1 greater than remaining buf len",
		},
		"string/1, not enough bytes": {
			[]byte{3<<5 | 24, 1},
			"slice len 1 greater than remaining buf len",
		},
		"string/?, no break": {
			[]byte{3<<5 | 31},
			"expected break marker",
		},
		"string/?, invalid nested major": {
			[]byte{3<<5 | 31, 2<<5 | 0},
			"unexpected major type 2 in indefinite slice",
		},
		"string/?, nested indefinite": {
			[]byte{3<<5 | 31, 3<<5 | 31},
			"nested indefinite slice",
		},
		"string/?, invalid nested definite": {
			[]byte{3<<5 | 31, 3<<5 | 24, 1},
			"decode subslice: slice len 1 greater than remaining buf len",
		},
	} {
		t.Run(name, func(t *testing.T) {
			cases = append(cases, errorTest{
				Description: fmt.Sprintf("TestDecode_InvalidSlice - %s - %s", name, c.Err),
				Input:       hex.EncodeToString(c.In),
			})

			_, _, err := decode(c.In)
			if err == nil {
				t.Errorf("expect err %s", c.Err)
			}
			if aerr := err.Error(); !strings.Contains(aerr, c.Err) {
				t.Errorf("expect err %s, got %s", c.Err, aerr)
			}
		})
	}

	dumpErrorCases(t, "TestDecodeError_InvalidSlice", cases)
}

func TestDecode_InvalidList(t *testing.T) {
	cases := []errorTest{}
	for name, c := range map[string]struct {
		In  []byte
		Err string
	}{
		"[] / eof after head": {
			[]byte{4<<5 | 1},
			"unexpected end of payload",
		},
		"[] / invalid item": {
			[]byte{4<<5 | 1, 0<<5 | 24},
			"arg len 1 greater than remaining buf len",
		},
		"[_ ] / no break": {
			[]byte{4<<5 | 31},
			"expected break marker",
		},
		"[_ ] / invalid item": {
			[]byte{4<<5 | 31, 0<<5 | 24},
			"arg len 1 greater than remaining buf len",
		},
	} {
		t.Run(name, func(t *testing.T) {
			cases = append(cases, errorTest{
				Description: fmt.Sprintf("TestDecode_InvalidList - %s - %s", name, c.Err),
				Input:       hex.EncodeToString(c.In),
			})

			_, _, err := decode(c.In)
			if err == nil {
				t.Errorf("expect err %s", c.Err)
			}
			if aerr := err.Error(); !strings.Contains(aerr, c.Err) {
				t.Errorf("expect err %s, got %s", c.Err, aerr)
			}
		})
	}

	dumpErrorCases(t, "TestDecodeError_InvalidList", cases)
}

func TestDecode_InvalidMap(t *testing.T) {
	cases := []errorTest{}
	for name, c := range map[string]struct {
		In  []byte
		Err string
	}{
		"{} / eof after head": {
			[]byte{5<<5 | 1},
			"unexpected end of payload",
		},
		"{} / non-string key": {
			[]byte{5<<5 | 1, 0},
			"unexpected major type 0 for map key",
		},
		"{} / invalid key": {
			[]byte{5<<5 | 1, 3<<5 | 24, 1},
			"slice len 1 greater than remaining buf len",
		},
		"{} / invalid value": {
			[]byte{5<<5 | 1, 3<<5 | 3, 0x66, 0x6f, 0x6f, 0<<5 | 24},
			"arg len 1 greater than remaining buf len",
		},
		"{_ } / no break": {
			[]byte{5<<5 | 31},
			"expected break marker",
		},
		"{_ } / non-string key": {
			[]byte{5<<5 | 31, 0},
			"unexpected major type 0 for map key",
		},
		"{_ } / invalid key": {
			[]byte{5<<5 | 31, 3<<5 | 24, 1},
			"slice len 1 greater than remaining buf len",
		},
		"{_ } / invalid value": {
			[]byte{5<<5 | 31, 3<<5 | 3, 0x66, 0x6f, 0x6f, 0<<5 | 24},
			"arg len 1 greater than remaining buf len",
		},
	} {
		t.Run(name, func(t *testing.T) {
			cases = append(cases, errorTest{
				Description: fmt.Sprintf("TestDecode_InvalidMap - %s - %s", name, c.Err),
				Input:       hex.EncodeToString(c.In),
			})

			_, _, err := decode(c.In)
			if err == nil {
				t.Errorf("expect err %s", c.Err)
			}
			if aerr := err.Error(); !strings.Contains(aerr, c.Err) {
				t.Errorf("expect err %s, got %s", c.Err, aerr)
			}
		})
	}

	dumpErrorCases(t, "TestDecodeError_InvalidMap", cases)
}

func TestDecode_InvalidTag(t *testing.T) {
	cases := []errorTest{}
	for name, c := range map[string]struct {
		In  []byte
		Err string
	}{
		"invalid value": {
			[]byte{6<<5 | 1, 0<<5 | 24},
			"arg len 1 greater than remaining buf len",
		},
		"eof": {
			[]byte{6<<5 | 1},
			"unexpected end of payload",
		},
	} {
		t.Run(name, func(t *testing.T) {
			cases = append(cases, errorTest{
				Description: fmt.Sprintf("TestDecode_InvalidTag - %s - %s", name, c.Err),
				Input:       hex.EncodeToString(c.In),
			})

			_, _, err := decode(c.In)
			if err == nil {
				t.Errorf("expect err %s", c.Err)
			}
			if aerr := err.Error(); !strings.Contains(aerr, c.Err) {
				t.Errorf("expect err %s, got %s", c.Err, aerr)
			}
		})
	}

	dumpErrorCases(t, "TestDecodeError_InvalidTag", cases)
}

func TestDecode_Atomic(t *testing.T) {
	cases := []successTest{}
	for name, c := range map[string]struct {
		In     []byte
		Expect Value
	}{
		"uint/0/min": {
			[]byte{0<<5 | 0},
			Uint(0),
		},
		"uint/0/max": {
			[]byte{0<<5 | 23},
			Uint(23),
		},
		"uint/1/min": {
			[]byte{0<<5 | 24, 0},
			Uint(0),
		},
		"uint/1/max": {
			[]byte{0<<5 | 24, 0xff},
			Uint(0xff),
		},
		"uint/2/min": {
			[]byte{0<<5 | 25, 0, 0},
			Uint(0),
		},
		"uint/2/max": {
			[]byte{0<<5 | 25, 0xff, 0xff},
			Uint(0xffff),
		},
		"uint/4/min": {
			[]byte{0<<5 | 26, 0, 0, 0, 0},
			Uint(0),
		},
		"uint/4/max": {
			[]byte{0<<5 | 26, 0xff, 0xff, 0xff, 0xff},
			Uint(0xffffffff),
		},
		"uint/8/min": {
			[]byte{0<<5 | 27, 0, 0, 0, 0, 0, 0, 0, 0},
			Uint(0),
		},
		"uint/8/max": {
			[]byte{0<<5 | 27, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			Uint(0xffffffff_ffffffff),
		},
		"negint/0/min": {
			[]byte{1<<5 | 0},
			NegInt(1),
		},
		"negint/0/max": {
			[]byte{1<<5 | 23},
			NegInt(24),
		},
		"negint/1/min": {
			[]byte{1<<5 | 24, 0},
			NegInt(1),
		},
		"negint/1/max": {
			[]byte{1<<5 | 24, 0xff},
			NegInt(0x100),
		},
		"negint/2/min": {
			[]byte{1<<5 | 25, 0, 0},
			NegInt(1),
		},
		"negint/2/max": {
			[]byte{1<<5 | 25, 0xff, 0xff},
			NegInt(0x10000),
		},
		"negint/4/min": {
			[]byte{1<<5 | 26, 0, 0, 0, 0},
			NegInt(1),
		},
		"negint/4/max": {
			[]byte{1<<5 | 26, 0xff, 0xff, 0xff, 0xff},
			NegInt(0x100000000),
		},
		"negint/8/min": {
			[]byte{1<<5 | 27, 0, 0, 0, 0, 0, 0, 0, 0},
			NegInt(1),
		},
		"negint/8/max": {
			[]byte{1<<5 | 27, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe},
			NegInt(0xffffffff_ffffffff),
		},
		"true": {
			[]byte{7<<5 | major7True},
			Bool(true),
		},
		"false": {
			[]byte{7<<5 | major7False},
			Bool(false),
		},
		"null": {
			[]byte{7<<5 | major7Nil},
			&Nil{},
		},
		"undefined": {
			[]byte{7<<5 | major7Undefined},
			&Undefined{},
		},
		"float16/+Inf": {
			[]byte{7<<5 | major7Float16, 0x7c, 0},
			Float32(math.Float32frombits(0x7f800000)),
		},
		"float16/-Inf": {
			[]byte{7<<5 | major7Float16, 0xfc, 0},
			Float32(math.Float32frombits(0xff800000)),
		},
		"float16/NaN/MSB": {
			[]byte{7<<5 | major7Float16, 0x7e, 0},
			Float32(math.Float32frombits(0x7fc00000)),
		},
		"float16/NaN/LSB": {
			[]byte{7<<5 | major7Float16, 0x7c, 1},
			Float32(math.Float32frombits(0x7f802000)),
		},
		"float32": {
			[]byte{7<<5 | major7Float32, 0x7f, 0x80, 0, 0},
			Float32(math.Float32frombits(0x7f800000)),
		},
		"float64": {
			[]byte{7<<5 | major7Float64, 0x7f, 0xf0, 0, 0, 0, 0, 0, 0},
			Float64(math.Float64frombits(0x7ff00000_00000000)),
		},
	} {
		t.Run(name, func(t *testing.T) {
			cases = append(cases, successTest{
				Description: fmt.Sprintf("atomic - %s", name),
				Input:       hex.EncodeToString(c.In),
				Expect:      toExpect(c.Expect),
			})

			actual, n, err := decode(c.In)
			if err != nil {
				t.Errorf("expect no err, got %v", err)
			}
			if n != len(c.In) {
				t.Errorf("didn't decode whole buffer")
			}
			assertValue(t, c.Expect, actual)
		})
	}

	dumpCases(t, "TestDecode_Atomic", cases)
}

func TestDecode_DefiniteSlice(t *testing.T) {
	cases := []successTest{}
	for name, c := range map[string]struct {
		In     []byte
		Expect Value
	}{
		"len = 0": {
			[]byte{2<<5 | 0},
			Slice{},
		},
		"len > 0": {
			[]byte{2<<5 | 3, 0x66, 0x6f, 0x6f},
			Slice{0x66, 0x6f, 0x6f},
		},
	} {
		t.Run(name, func(t *testing.T) {
			cases = append(cases, successTest{
				Description: fmt.Sprintf("definite slice - %s", name),
				Input:       hex.EncodeToString(c.In),
				Expect:      toExpect(c.Expect),
			})

			actual, n, err := decode(c.In)
			if err != nil {
				t.Errorf("expect no err, got %v", err)
			}
			if n != len(c.In) {
				t.Errorf("didn't decode whole buffer")
			}
			assertValue(t, c.Expect, actual)
		})
	}

	dumpCases(t, "TestDecode_DefiniteSlice", cases)
}

func TestDecode_IndefiniteSlice(t *testing.T) {
	cases := []successTest{}
	for name, c := range map[string]struct {
		In     []byte
		Expect Value
	}{
		"len = 0": {
			[]byte{2<<5 | 31, 0xff},
			Slice{},
		},
		"len = 0, explicit": {
			[]byte{2<<5 | 31, 2<<5 | 0, 0xff},
			Slice{},
		},
		"len = 0, len > 0": {
			[]byte{
				2<<5 | 31,
				2<<5 | 0,
				2<<5 | 3, 0x66, 0x6f, 0x6f,
				0xff,
			},
			Slice{0x66, 0x6f, 0x6f},
		},
		"len > 0, len = 0": {
			[]byte{
				2<<5 | 31,
				2<<5 | 3, 0x66, 0x6f, 0x6f,
				2<<5 | 0,
				0xff,
			},
			Slice{0x66, 0x6f, 0x6f},
		},
		"len > 0, len > 0": {
			[]byte{
				2<<5 | 31,
				2<<5 | 3, 0x66, 0x6f, 0x6f,
				2<<5 | 3, 0x66, 0x6f, 0x6f,
				0xff,
			},
			Slice{0x66, 0x6f, 0x6f, 0x66, 0x6f, 0x6f},
		},
	} {
		t.Run(name, func(t *testing.T) {
			cases = append(cases, successTest{
				Description: fmt.Sprintf("indefinite slice - %s", name),
				Input:       hex.EncodeToString(c.In),
				Expect:      toExpect(c.Expect),
			})

			actual, n, err := decode(c.In)
			if err != nil {
				t.Errorf("expect no err, got %v", err)
			}
			if n != len(c.In) {
				t.Errorf("didn't decode whole buffer")
			}
			assertValue(t, c.Expect, actual)
		})
	}

	dumpCases(t, "TestDecode_IndefiniteSlice", cases)
}

func TestDecode_DefiniteString(t *testing.T) {
	cases := []successTest{}
	for name, c := range map[string]struct {
		In     []byte
		Expect Value
	}{
		"len = 0": {
			[]byte{3<<5 | 0},
			String(""),
		},
		"len > 0": {
			[]byte{3<<5 | 3, 0x66, 0x6f, 0x6f},
			String("foo"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			cases = append(cases, successTest{
				Description: fmt.Sprintf("definite string - %s", name),
				Input:       hex.EncodeToString(c.In),
				Expect:      toExpect(c.Expect),
			})

			actual, n, err := decode(c.In)
			if err != nil {
				t.Errorf("expect no err, got %v", err)
			}
			if n != len(c.In) {
				t.Errorf("didn't decode whole buffer")
			}
			assertValue(t, c.Expect, actual)
		})
	}

	dumpCases(t, "TestDecode_DefiniteString", cases)
}

func TestDecode_IndefiniteString(t *testing.T) {
	cases := []successTest{}
	for name, c := range map[string]struct {
		In     []byte
		Expect Value
	}{
		"len = 0": {
			[]byte{3<<5 | 31, 0xff},
			String(""),
		},
		"len = 0, explicit": {
			[]byte{3<<5 | 31, 3<<5 | 0, 0xff},
			String(""),
		},
		"len = 0, len > 0": {
			[]byte{
				3<<5 | 31,
				3<<5 | 0,
				3<<5 | 3, 0x66, 0x6f, 0x6f,
				0xff,
			},
			String("foo"),
		},
		"len > 0, len = 0": {
			[]byte{
				3<<5 | 31,
				3<<5 | 3, 0x66, 0x6f, 0x6f,
				3<<5 | 0,
				0xff,
			},
			String("foo"),
		},
		"len > 0, len > 0": {
			[]byte{
				3<<5 | 31,
				3<<5 | 3, 0x66, 0x6f, 0x6f,
				3<<5 | 3, 0x66, 0x6f, 0x6f,
				0xff,
			},
			String("foofoo"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			cases = append(cases, successTest{
				Description: fmt.Sprintf("indefinite string - %s", name),
				Input:       hex.EncodeToString(c.In),
				Expect:      toExpect(c.Expect),
			})

			actual, n, err := decode(c.In)
			if err != nil {
				t.Errorf("expect no err, got %v", err)
			}
			if n != len(c.In) {
				t.Errorf("didn't decode whole buffer")
			}
			assertValue(t, c.Expect, actual)
		})
	}

	dumpCases(t, "TestDecode_IndefiniteString", cases)
}

func TestDecode_List(t *testing.T) {
	cases := []successTest{}
	for name, c := range map[string]struct {
		In     []byte
		Expect Value
	}{
		"[uint/0/min]": {
			In:     withDefiniteList([]byte{0<<5 | 0}),
			Expect: List{Uint(0)},
		},
		"[uint/0/max]": {
			In:     withDefiniteList([]byte{0<<5 | 23}),
			Expect: List{Uint(23)},
		},
		"[uint/1/min]": {
			In:     withDefiniteList([]byte{0<<5 | 24, 0}),
			Expect: List{Uint(0)},
		},
		"[uint/1/max]": {
			In:     withDefiniteList([]byte{0<<5 | 24, 0xff}),
			Expect: List{Uint(0xff)},
		},
		"[uint/2/min]": {
			In:     withDefiniteList([]byte{0<<5 | 25, 0, 0}),
			Expect: List{Uint(0)},
		},
		"[uint/2/max]": {
			In:     withDefiniteList([]byte{0<<5 | 25, 0xff, 0xff}),
			Expect: List{Uint(0xffff)},
		},
		"[uint/4/min]": {
			In:     withDefiniteList([]byte{0<<5 | 26, 0, 0, 0, 0}),
			Expect: List{Uint(0)},
		},
		"[uint/4/max]": {
			In:     withDefiniteList([]byte{0<<5 | 26, 0xff, 0xff, 0xff, 0xff}),
			Expect: List{Uint(0xffffffff)},
		},
		"[uint/8/min]": {
			In:     withDefiniteList([]byte{0<<5 | 27, 0, 0, 0, 0, 0, 0, 0, 0}),
			Expect: List{Uint(0)},
		},
		"[uint/8/max]": {
			In:     withDefiniteList([]byte{0<<5 | 27, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}),
			Expect: List{Uint(0xffffffff_ffffffff)},
		},
		"[negint/0/min]": {
			In:     withDefiniteList([]byte{1<<5 | 0}),
			Expect: List{NegInt(1)},
		},
		"[negint/0/max]": {
			In:     withDefiniteList([]byte{1<<5 | 23}),
			Expect: List{NegInt(24)},
		},
		"[negint/1/min]": {
			In:     withDefiniteList([]byte{1<<5 | 24, 0}),
			Expect: List{NegInt(1)},
		},
		"[negint/1/max]": {
			In:     withDefiniteList([]byte{1<<5 | 24, 0xff}),
			Expect: List{NegInt(0x100)},
		},
		"[negint/2/min]": {
			In:     withDefiniteList([]byte{1<<5 | 25, 0, 0}),
			Expect: List{NegInt(1)},
		},
		"[negint/2/max]": {
			In:     withDefiniteList([]byte{1<<5 | 25, 0xff, 0xff}),
			Expect: List{NegInt(0x10000)},
		},
		"[negint/4/min]": {
			In:     withDefiniteList([]byte{1<<5 | 26, 0, 0, 0, 0}),
			Expect: List{NegInt(1)},
		},
		"[negint/4/max]": {
			In:     withDefiniteList([]byte{1<<5 | 26, 0xff, 0xff, 0xff, 0xff}),
			Expect: List{NegInt(0x100000000)},
		},
		"[negint/8/min]": {
			In:     withDefiniteList([]byte{1<<5 | 27, 0, 0, 0, 0, 0, 0, 0, 0}),
			Expect: List{NegInt(1)},
		},
		"[negint/8/max]": {
			In:     withDefiniteList([]byte{1<<5 | 27, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe}),
			Expect: List{NegInt(0xffffffff_ffffffff)},
		},
		"[true]": {
			In:     withDefiniteList([]byte{7<<5 | major7True}),
			Expect: List{Bool(true)},
		},
		"[false]": {
			In:     withDefiniteList([]byte{7<<5 | major7False}),
			Expect: List{Bool(false)},
		},
		"[null]": {
			In:     withDefiniteList([]byte{7<<5 | major7Nil}),
			Expect: List{&Nil{}},
		},
		"[undefined]": {
			In:     withDefiniteList([]byte{7<<5 | major7Undefined}),
			Expect: List{&Undefined{}},
		},
		"[float16/+Inf]": {
			In:     withDefiniteList([]byte{7<<5 | major7Float16, 0x7c, 0}),
			Expect: List{Float32(math.Float32frombits(0x7f800000))},
		},
		"[float16/-Inf]": {
			In:     withDefiniteList([]byte{7<<5 | major7Float16, 0xfc, 0}),
			Expect: List{Float32(math.Float32frombits(0xff800000))},
		},
		"[float16/NaN/MSB]": {
			In:     withDefiniteList([]byte{7<<5 | major7Float16, 0x7e, 0}),
			Expect: List{Float32(math.Float32frombits(0x7fc00000))},
		},
		"[float16/NaN/LSB]": {
			In:     withDefiniteList([]byte{7<<5 | major7Float16, 0x7c, 1}),
			Expect: List{Float32(math.Float32frombits(0x7f802000))},
		},
		"[float32]": {
			In:     withDefiniteList([]byte{7<<5 | major7Float32, 0x7f, 0x80, 0, 0}),
			Expect: List{Float32(math.Float32frombits(0x7f800000))},
		},
		"[float64]": {
			In:     withDefiniteList([]byte{7<<5 | major7Float64, 0x7f, 0xf0, 0, 0, 0, 0, 0, 0}),
			Expect: List{Float64(math.Float64frombits(0x7ff00000_00000000))},
		},
		"[_ uint/0/min]": {
			In:     withIndefiniteList([]byte{0<<5 | 0}),
			Expect: List{Uint(0)},
		},
		"[_ uint/0/max]": {
			In:     withIndefiniteList([]byte{0<<5 | 23}),
			Expect: List{Uint(23)},
		},
		"[_ uint/1/min]": {
			In:     withIndefiniteList([]byte{0<<5 | 24, 0}),
			Expect: List{Uint(0)},
		},
		"[_ uint/1/max]": {
			In:     withIndefiniteList([]byte{0<<5 | 24, 0xff}),
			Expect: List{Uint(0xff)},
		},
		"[_ uint/2/min]": {
			In:     withIndefiniteList([]byte{0<<5 | 25, 0, 0}),
			Expect: List{Uint(0)},
		},
		"[_ uint/2/max]": {
			In:     withIndefiniteList([]byte{0<<5 | 25, 0xff, 0xff}),
			Expect: List{Uint(0xffff)},
		},
		"[_ uint/4/min]": {
			In:     withIndefiniteList([]byte{0<<5 | 26, 0, 0, 0, 0}),
			Expect: List{Uint(0)},
		},
		"[_ uint/4/max]": {
			In:     withIndefiniteList([]byte{0<<5 | 26, 0xff, 0xff, 0xff, 0xff}),
			Expect: List{Uint(0xffffffff)},
		},
		"[_ uint/8/min]": {
			In:     withIndefiniteList([]byte{0<<5 | 27, 0, 0, 0, 0, 0, 0, 0, 0}),
			Expect: List{Uint(0)},
		},
		"[_ uint/8/max]": {
			In:     withIndefiniteList([]byte{0<<5 | 27, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}),
			Expect: List{Uint(0xffffffff_ffffffff)},
		},
		"[_ negint/0/min]": {
			In:     withIndefiniteList([]byte{1<<5 | 0}),
			Expect: List{NegInt(1)},
		},
		"[_ negint/0/max]": {
			In:     withIndefiniteList([]byte{1<<5 | 23}),
			Expect: List{NegInt(24)},
		},
		"[_ negint/1/min]": {
			In:     withIndefiniteList([]byte{1<<5 | 24, 0}),
			Expect: List{NegInt(1)},
		},
		"[_ negint/1/max]": {
			In:     withIndefiniteList([]byte{1<<5 | 24, 0xff}),
			Expect: List{NegInt(0x100)},
		},
		"[_ negint/2/min]": {
			In:     withIndefiniteList([]byte{1<<5 | 25, 0, 0}),
			Expect: List{NegInt(1)},
		},
		"[_ negint/2/max]": {
			In:     withIndefiniteList([]byte{1<<5 | 25, 0xff, 0xff}),
			Expect: List{NegInt(0x10000)},
		},
		"[_ negint/4/min]": {
			In:     withIndefiniteList([]byte{1<<5 | 26, 0, 0, 0, 0}),
			Expect: List{NegInt(1)},
		},
		"[_ negint/4/max]": {
			In:     withIndefiniteList([]byte{1<<5 | 26, 0xff, 0xff, 0xff, 0xff}),
			Expect: List{NegInt(0x100000000)},
		},
		"[_ negint/8/min]": {
			In:     withIndefiniteList([]byte{1<<5 | 27, 0, 0, 0, 0, 0, 0, 0, 0}),
			Expect: List{NegInt(1)},
		},
		"[_ negint/8/max]": {
			In:     withIndefiniteList([]byte{1<<5 | 27, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe}),
			Expect: List{NegInt(0xffffffff_ffffffff)},
		},
		"[_ true]": {
			In:     withIndefiniteList([]byte{7<<5 | major7True}),
			Expect: List{Bool(true)},
		},
		"[_ false]": {
			In:     withIndefiniteList([]byte{7<<5 | major7False}),
			Expect: List{Bool(false)},
		},
		"[_ null]": {
			In:     withIndefiniteList([]byte{7<<5 | major7Nil}),
			Expect: List{&Nil{}},
		},
		"[_ undefined]": {
			In:     withIndefiniteList([]byte{7<<5 | major7Undefined}),
			Expect: List{&Undefined{}},
		},
		"[_ float16/+Inf]": {
			In:     withIndefiniteList([]byte{7<<5 | major7Float16, 0x7c, 0}),
			Expect: List{Float32(math.Float32frombits(0x7f800000))},
		},
		"[_ float16/-Inf]": {
			In:     withIndefiniteList([]byte{7<<5 | major7Float16, 0xfc, 0}),
			Expect: List{Float32(math.Float32frombits(0xff800000))},
		},
		"[_ float16/NaN/MSB]": {
			In:     withIndefiniteList([]byte{7<<5 | major7Float16, 0x7e, 0}),
			Expect: List{Float32(math.Float32frombits(0x7fc00000))},
		},
		"[_ float16/NaN/LSB]": {
			In:     withIndefiniteList([]byte{7<<5 | major7Float16, 0x7c, 1}),
			Expect: List{Float32(math.Float32frombits(0x7f802000))},
		},
		"[_ float32]": {
			In:     withIndefiniteList([]byte{7<<5 | major7Float32, 0x7f, 0x80, 0, 0}),
			Expect: List{Float32(math.Float32frombits(0x7f800000))},
		},
		"[_ float64]": {
			In:     withIndefiniteList([]byte{7<<5 | major7Float64, 0x7f, 0xf0, 0, 0, 0, 0, 0, 0}),
			Expect: List{Float64(math.Float64frombits(0x7ff00000_00000000))},
		},
	} {
		t.Run(name, func(t *testing.T) {
			cases = append(cases, successTest{
				Description: fmt.Sprintf("list - %s", name),
				Input:       hex.EncodeToString(c.In),
				Expect:      toExpect(c.Expect),
			})

			actual, n, err := decode(c.In)
			if err != nil {
				t.Errorf("expect no err, got %v", err)
			}
			if n != len(c.In) {
				t.Errorf("didn't decode whole buffer (decoded %d of %d)", n, len(c.In))
			}
			assertValue(t, c.Expect, actual)
		})
	}

	dumpCases(t, "TestDecode_List", cases)
}

func TestDecode_Map(t *testing.T) {
	cases := []successTest{}
	for name, c := range map[string]struct {
		In     []byte
		Expect Value
	}{
		"{uint/0/min}": {
			In:     withDefiniteMap([]byte{0<<5 | 0}),
			Expect: Map{"foo": Uint(0)},
		},
		"{uint/0/max}": {
			In:     withDefiniteMap([]byte{0<<5 | 23}),
			Expect: Map{"foo": Uint(23)},
		},
		"{uint/1/min}": {
			In:     withDefiniteMap([]byte{0<<5 | 24, 0}),
			Expect: Map{"foo": Uint(0)},
		},
		"{uint/1/max}": {
			In:     withDefiniteMap([]byte{0<<5 | 24, 0xff}),
			Expect: Map{"foo": Uint(0xff)},
		},
		"{uint/2/min}": {
			In:     withDefiniteMap([]byte{0<<5 | 25, 0, 0}),
			Expect: Map{"foo": Uint(0)},
		},
		"{uint/2/max}": {
			In:     withDefiniteMap([]byte{0<<5 | 25, 0xff, 0xff}),
			Expect: Map{"foo": Uint(0xffff)},
		},
		"{uint/4/min}": {
			In:     withDefiniteMap([]byte{0<<5 | 26, 0, 0, 0, 0}),
			Expect: Map{"foo": Uint(0)},
		},
		"{uint/4/max}": {
			In:     withDefiniteMap([]byte{0<<5 | 26, 0xff, 0xff, 0xff, 0xff}),
			Expect: Map{"foo": Uint(0xffffffff)},
		},
		"{uint/8/min}": {
			In:     withDefiniteMap([]byte{0<<5 | 27, 0, 0, 0, 0, 0, 0, 0, 0}),
			Expect: Map{"foo": Uint(0)},
		},
		"{uint/8/max}": {
			In:     withDefiniteMap([]byte{0<<5 | 27, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}),
			Expect: Map{"foo": Uint(0xffffffff_ffffffff)},
		},
		"{negint/0/min}": {
			In:     withDefiniteMap([]byte{1<<5 | 0}),
			Expect: Map{"foo": NegInt(1)},
		},
		"{negint/0/max}": {
			In:     withDefiniteMap([]byte{1<<5 | 23}),
			Expect: Map{"foo": NegInt(24)},
		},
		"{negint/1/min}": {
			In:     withDefiniteMap([]byte{1<<5 | 24, 0}),
			Expect: Map{"foo": NegInt(1)},
		},
		"{negint/1/max}": {
			In:     withDefiniteMap([]byte{1<<5 | 24, 0xff}),
			Expect: Map{"foo": NegInt(0x100)},
		},
		"{negint/2/min}": {
			In:     withDefiniteMap([]byte{1<<5 | 25, 0, 0}),
			Expect: Map{"foo": NegInt(1)},
		},
		"{negint/2/max}": {
			In:     withDefiniteMap([]byte{1<<5 | 25, 0xff, 0xff}),
			Expect: Map{"foo": NegInt(0x10000)},
		},
		"{negint/4/min}": {
			In:     withDefiniteMap([]byte{1<<5 | 26, 0, 0, 0, 0}),
			Expect: Map{"foo": NegInt(1)},
		},
		"{negint/4/max}": {
			In:     withDefiniteMap([]byte{1<<5 | 26, 0xff, 0xff, 0xff, 0xff}),
			Expect: Map{"foo": NegInt(0x100000000)},
		},
		"{negint/8/min}": {
			In:     withDefiniteMap([]byte{1<<5 | 27, 0, 0, 0, 0, 0, 0, 0, 0}),
			Expect: Map{"foo": NegInt(1)},
		},
		"{negint/8/max}": {
			In:     withDefiniteMap([]byte{1<<5 | 27, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe}),
			Expect: Map{"foo": NegInt(0xffffffff_ffffffff)},
		},
		"{true}": {
			In:     withDefiniteMap([]byte{7<<5 | major7True}),
			Expect: Map{"foo": Bool(true)},
		},
		"{false}": {
			In:     withDefiniteMap([]byte{7<<5 | major7False}),
			Expect: Map{"foo": Bool(false)},
		},
		"{null}": {
			In:     withDefiniteMap([]byte{7<<5 | major7Nil}),
			Expect: Map{"foo": &Nil{}},
		},
		"{undefined}": {
			In:     withDefiniteMap([]byte{7<<5 | major7Undefined}),
			Expect: Map{"foo": &Undefined{}},
		},
		"{float16/+Inf}": {
			In:     withDefiniteMap([]byte{7<<5 | major7Float16, 0x7c, 0}),
			Expect: Map{"foo": Float32(math.Float32frombits(0x7f800000))},
		},
		"{float16/-Inf}": {
			In:     withDefiniteMap([]byte{7<<5 | major7Float16, 0xfc, 0}),
			Expect: Map{"foo": Float32(math.Float32frombits(0xff800000))},
		},
		"{float16/NaN/MSB}": {
			In:     withDefiniteMap([]byte{7<<5 | major7Float16, 0x7e, 0}),
			Expect: Map{"foo": Float32(math.Float32frombits(0x7fc00000))},
		},
		"{float16/NaN/LSB}": {
			In:     withDefiniteMap([]byte{7<<5 | major7Float16, 0x7c, 1}),
			Expect: Map{"foo": Float32(math.Float32frombits(0x7f802000))},
		},
		"{float32}": {
			In:     withDefiniteMap([]byte{7<<5 | major7Float32, 0x7f, 0x80, 0, 0}),
			Expect: Map{"foo": Float32(math.Float32frombits(0x7f800000))},
		},
		"{float64}": {
			In:     withDefiniteMap([]byte{7<<5 | major7Float64, 0x7f, 0xf0, 0, 0, 0, 0, 0, 0}),
			Expect: Map{"foo": Float64(math.Float64frombits(0x7ff00000_00000000))},
		},
		"{_ uint/0/min}": {
			In:     withIndefiniteMap([]byte{0<<5 | 0}),
			Expect: Map{"foo": Uint(0)},
		},
		"{_ uint/0/max}": {
			In:     withIndefiniteMap([]byte{0<<5 | 23}),
			Expect: Map{"foo": Uint(23)},
		},
		"{_ uint/1/min}": {
			In:     withIndefiniteMap([]byte{0<<5 | 24, 0}),
			Expect: Map{"foo": Uint(0)},
		},
		"{_ uint/1/max}": {
			In:     withIndefiniteMap([]byte{0<<5 | 24, 0xff}),
			Expect: Map{"foo": Uint(0xff)},
		},
		"{_ uint/2/min}": {
			In:     withIndefiniteMap([]byte{0<<5 | 25, 0, 0}),
			Expect: Map{"foo": Uint(0)},
		},
		"{_ uint/2/max}": {
			In:     withIndefiniteMap([]byte{0<<5 | 25, 0xff, 0xff}),
			Expect: Map{"foo": Uint(0xffff)},
		},
		"{_ uint/4/min}": {
			In:     withIndefiniteMap([]byte{0<<5 | 26, 0, 0, 0, 0}),
			Expect: Map{"foo": Uint(0)},
		},
		"{_ uint/4/max}": {
			In:     withIndefiniteMap([]byte{0<<5 | 26, 0xff, 0xff, 0xff, 0xff}),
			Expect: Map{"foo": Uint(0xffffffff)},
		},
		"{_ uint/8/min}": {
			In:     withIndefiniteMap([]byte{0<<5 | 27, 0, 0, 0, 0, 0, 0, 0, 0}),
			Expect: Map{"foo": Uint(0)},
		},
		"{_ uint/8/max}": {
			In:     withIndefiniteMap([]byte{0<<5 | 27, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}),
			Expect: Map{"foo": Uint(0xffffffff_ffffffff)},
		},
		"{_ negint/0/min}": {
			In:     withIndefiniteMap([]byte{1<<5 | 0}),
			Expect: Map{"foo": NegInt(1)},
		},
		"{_ negint/0/max}": {
			In:     withIndefiniteMap([]byte{1<<5 | 23}),
			Expect: Map{"foo": NegInt(24)},
		},
		"{_ negint/1/min}": {
			In:     withIndefiniteMap([]byte{1<<5 | 24, 0}),
			Expect: Map{"foo": NegInt(1)},
		},
		"{_ negint/1/max}": {
			In:     withIndefiniteMap([]byte{1<<5 | 24, 0xff}),
			Expect: Map{"foo": NegInt(0x100)},
		},
		"{_ negint/2/min}": {
			In:     withIndefiniteMap([]byte{1<<5 | 25, 0, 0}),
			Expect: Map{"foo": NegInt(1)},
		},
		"{_ negint/2/max}": {
			In:     withIndefiniteMap([]byte{1<<5 | 25, 0xff, 0xff}),
			Expect: Map{"foo": NegInt(0x10000)},
		},
		"{_ negint/4/min}": {
			In:     withIndefiniteMap([]byte{1<<5 | 26, 0, 0, 0, 0}),
			Expect: Map{"foo": NegInt(1)},
		},
		"{_ negint/4/max}": {
			In:     withIndefiniteMap([]byte{1<<5 | 26, 0xff, 0xff, 0xff, 0xff}),
			Expect: Map{"foo": NegInt(0x100000000)},
		},
		"{_ negint/8/min}": {
			In:     withIndefiniteMap([]byte{1<<5 | 27, 0, 0, 0, 0, 0, 0, 0, 0}),
			Expect: Map{"foo": NegInt(1)},
		},
		"{_ negint/8/max}": {
			In:     withIndefiniteMap([]byte{1<<5 | 27, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe}),
			Expect: Map{"foo": NegInt(0xffffffff_ffffffff)},
		},
		"{_ true}": {
			In:     withIndefiniteMap([]byte{7<<5 | major7True}),
			Expect: Map{"foo": Bool(true)},
		},
		"{_ false}": {
			In:     withIndefiniteMap([]byte{7<<5 | major7False}),
			Expect: Map{"foo": Bool(false)},
		},
		"{_ null}": {
			In:     withIndefiniteMap([]byte{7<<5 | major7Nil}),
			Expect: Map{"foo": &Nil{}},
		},
		"{_ undefined}": {
			In:     withIndefiniteMap([]byte{7<<5 | major7Undefined}),
			Expect: Map{"foo": &Undefined{}},
		},
		"{_ float16/+Inf}": {
			In:     withIndefiniteMap([]byte{7<<5 | major7Float16, 0x7c, 0}),
			Expect: Map{"foo": Float32(math.Float32frombits(0x7f800000))},
		},
		"{_ float16/-Inf}": {
			In:     withIndefiniteMap([]byte{7<<5 | major7Float16, 0xfc, 0}),
			Expect: Map{"foo": Float32(math.Float32frombits(0xff800000))},
		},
		"{_ float16/NaN/MSB}": {
			In:     withIndefiniteMap([]byte{7<<5 | major7Float16, 0x7e, 0}),
			Expect: Map{"foo": Float32(math.Float32frombits(0x7fc00000))},
		},
		"{_ float16/NaN/LSB}": {
			In:     withIndefiniteMap([]byte{7<<5 | major7Float16, 0x7c, 1}),
			Expect: Map{"foo": Float32(math.Float32frombits(0x7f802000))},
		},
		"{_ float32}": {
			In:     withIndefiniteMap([]byte{7<<5 | major7Float32, 0x7f, 0x80, 0, 0}),
			Expect: Map{"foo": Float32(math.Float32frombits(0x7f800000))},
		},
		"{_ float64}": {
			In:     withIndefiniteMap([]byte{7<<5 | major7Float64, 0x7f, 0xf0, 0, 0, 0, 0, 0, 0}),
			Expect: Map{"foo": Float64(math.Float64frombits(0x7ff00000_00000000))},
		},
	} {
		t.Run(name, func(t *testing.T) {
			cases = append(cases, successTest{
				Description: fmt.Sprintf("map - %s", name),
				Input:       hex.EncodeToString(c.In),
				Expect:      toExpect(c.Expect),
			})

			actual, n, err := decode(c.In)
			if err != nil {
				t.Errorf("expect no err, got %v", err)
			}
			if n != len(c.In) {
				t.Errorf("didn't decode whole buffer (decoded %d of %d)", n, len(c.In))
			}
			assertValue(t, c.Expect, actual)
		})
	}

	dumpCases(t, "TestDecode_Map", cases)
}

func TestDecode_Tag(t *testing.T) {
	cases := []successTest{}
	for name, c := range map[string]struct {
		In     []byte
		Expect Value
	}{
		"0/min": {
			In:     []byte{6<<5 | 0, 1},
			Expect: &Tag{0, Uint(1)},
		},
		"0/max": {
			In:     []byte{6<<5 | 23, 1},
			Expect: &Tag{23, Uint(1)},
		},
		"1/min": {
			In:     []byte{6<<5 | 24, 0, 1},
			Expect: &Tag{0, Uint(1)},
		},
		"1/max": {
			In:     []byte{6<<5 | 24, 0xff, 1},
			Expect: &Tag{0xff, Uint(1)},
		},
		"2/min": {
			In:     []byte{6<<5 | 25, 0, 0, 1},
			Expect: &Tag{0, Uint(1)},
		},
		"2/max": {
			In:     []byte{6<<5 | 25, 0xff, 0xff, 1},
			Expect: &Tag{0xffff, Uint(1)},
		},
		"4/min": {
			In:     []byte{6<<5 | 26, 0, 0, 0, 0, 1},
			Expect: &Tag{0, Uint(1)},
		},
		"4/max": {
			In:     []byte{6<<5 | 26, 0xff, 0xff, 0xff, 0xff, 1},
			Expect: &Tag{0xffffffff, Uint(1)},
		},
		"8/min": {
			In:     []byte{6<<5 | 27, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			Expect: &Tag{0, Uint(1)},
		},
		"8/max": {
			In:     []byte{6<<5 | 27, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 1},
			Expect: &Tag{0xffffffff_ffffffff, Uint(1)},
		},
	} {
		t.Run(name, func(t *testing.T) {
			cases = append(cases, successTest{
				Description: fmt.Sprintf("tag - %s", name),
				Input:       hex.EncodeToString(c.In),
				Expect:      toExpect(c.Expect),
			})

			actual, n, err := decode(c.In)
			if err != nil {
				t.Errorf("expect no err, got %v", err)
			}
			if n != len(c.In) {
				t.Errorf("didn't decode whole buffer (decoded %d of %d)", n, len(c.In))
			}
			assertValue(t, c.Expect, actual)
		})
	}

	dumpCases(t, "TestDecode_Tag", cases)
}

func TestDecode_CollectSuccessTests(t *testing.T) {
	entries, err := os.ReadDir("./")
	if err != nil {
		t.Fatal(err)
	}

	cases := []successTest{}
	for _, e := range entries {
		ename := e.Name()
		if !strings.HasPrefix(ename, "TestDecode_") || !strings.HasSuffix(ename, ".json") {
			continue
		}

		f, err := os.Open(ename)
		if err != nil {
			t.Fatal(err)
		}

		p, err := io.ReadAll(f)
		if err != nil {
			t.Fatal(err)
		}

		var v []successTest
		if err := json.Unmarshal(p, &v); err != nil {
			t.Fatal(err)
		}

		fmt.Printf("append %d cases from %s\n", len(v), ename)
		cases = append(cases, v...)
	}

	dumpCases(t, "TestDecodeSuccess", cases)
}

func TestDecode_CollectErrorTests(t *testing.T) {
	entries, err := os.ReadDir("./")
	if err != nil {
		t.Fatal(err)
	}

	cases := []errorTest{}
	for _, e := range entries {
		ename := e.Name()
		if !strings.HasPrefix(ename, "TestDecodeError_") || !strings.HasSuffix(ename, ".json") {
			continue
		}

		f, err := os.Open(ename)
		if err != nil {
			t.Fatal(err)
		}

		p, err := io.ReadAll(f)
		if err != nil {
			t.Fatal(err)
		}

		var v []errorTest
		if err := json.Unmarshal(p, &v); err != nil {
			t.Fatal(err)
		}

		fmt.Printf("append %d cases from %s\n", len(v), ename)
		cases = append(cases, v...)
	}

	dumpErrorCases(t, "TestDecodeError", cases)
}

func dumpCases(t *testing.T, name string, cases []successTest) {
	f, err := os.Create(name + ".json")
	if err != nil {
		t.Fatal(err)
	}

	p, err := json.MarshalIndent(cases, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := f.Write(p); err != nil {
		t.Fatal(err)
	}
}

func dumpErrorCases(t *testing.T, name string, cases []errorTest) {
	f, err := os.Create(name + ".json")
	if err != nil {
		t.Fatal(err)
	}

	p, err := json.MarshalIndent(cases, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := f.Write(p); err != nil {
		t.Fatal(err)
	}
}

func assertValue(t *testing.T, e, a Value) {
	switch v := e.(type) {
	case Uint, NegInt, Slice, String, Bool, *Nil, *Undefined:
		if !reflect.DeepEqual(e, a) {
			t.Errorf("%v != %v", e, a)
		}
	case List:
		assertList(t, v, a)
	case Map:
		assertMap(t, v, a)
	case *Tag:
		assertTag(t, v, a)
	case Float32:
		assertMajor7Float32(t, v, a)
	case Float64:
		assertMajor7Float64(t, v, a)
	default:
		t.Errorf("unrecognized variant %T", e)
	}
}

func assertList(t *testing.T, e List, a Value) {
	av, ok := a.(List)
	if !ok {
		t.Errorf("%T != %T", e, a)
		return
	}

	if len(e) != len(av) {
		t.Errorf("length %d != %d", len(e), len(av))
		return
	}

	for i := 0; i < len(e); i++ {
		assertValue(t, e[i], av[i])
	}
}

func assertMap(t *testing.T, e Map, a Value) {
	av, ok := a.(Map)
	if !ok {
		t.Errorf("%T != %T", e, a)
		return
	}

	if len(e) != len(av) {
		t.Errorf("length %d != %d", len(e), len(av))
		return
	}

	for k, ev := range e {
		avv, ok := av[k]
		if !ok {
			t.Errorf("missing key %s", k)
			return
		}

		assertValue(t, ev, avv)
	}
}

func assertTag(t *testing.T, e *Tag, a Value) {
	av, ok := a.(*Tag)
	if !ok {
		t.Errorf("%T != %T", e, a)
		return
	}

	if e.ID != av.ID {
		t.Errorf("tag ID %d != %d", e.ID, av.ID)
		return
	}

	assertValue(t, e.Value, av.Value)
}

func assertMajor7Float32(t *testing.T, e Float32, a Value) {
	av, ok := a.(Float32)
	if !ok {
		t.Errorf("%T != %T", e, a)
		return
	}

	if math.Float32bits(float32(e)) != math.Float32bits(float32(av)) {
		t.Errorf("float32(%x) != float32(%x)", e, av)
	}
}

func assertMajor7Float64(t *testing.T, e Float64, a Value) {
	av, ok := a.(Float64)
	if !ok {
		t.Errorf("%T != %T", e, a)
		return
	}

	if math.Float64bits(float64(e)) != math.Float64bits(float64(av)) {
		t.Errorf("float64(%x) != float64(%x)", e, av)
	}
}

var mapKeyFoo = []byte{0x63, 0x66, 0x6f, 0x6f}

func withDefiniteList(p []byte) []byte {
	return append([]byte{4<<5 | 1}, p...)
}

func withIndefiniteList(p []byte) []byte {
	p = append([]byte{4<<5 | 31}, p...)
	return append(p, 0xff)
}

func withDefiniteMap(p []byte) []byte {
	head := append([]byte{5<<5 | 1}, mapKeyFoo...)
	return append(head, p...)
}

func withIndefiniteMap(p []byte) []byte {
	head := append([]byte{5<<5 | 31}, mapKeyFoo...)
	p = append(head, p...)
	return append(p, 0xff)
}

type errorTest struct {
	Description string `json:"description"`
	Input       string `json:"input"`
}

type successTest struct {
	Description string            `json:"description"`
	Input       string            `json:"input"`
	Expect      successTestExpect `json:"expect"`
}

type successTestExpect struct {
	// one of
	Uint       *expectUint      `json:"uint,omitempty"`
	Negint     *expectNegint    `json:"negint,omitempty"`
	ByteString expectByteString `json:"bytestring,omitempty"`
	String     *expectString    `json:"string,omitempty"`
	List       expectList       `json:"list,omitempty"`
	Map        expectMap        `json:"map,omitempty"`
	Tag        *expectTag       `json:"tag,omitempty"`
	Bool       *expectBool      `json:"bool,omitempty"`
	Null       *expectNull      `json:"null,omitempty"`
	Undefined  *expectUndefined `json:"undefined,omitempty"`
	Float32    *expectFloat32   `json:"float32,omitempty"`
	Float64    *expectFloat64   `json:"float64,omitempty"`
}

func toExpect(v Value) successTestExpect {
	expect := successTestExpect{}
	switch vv := v.(type) {
	case Uint:
		expect.Uint = (*expectUint)(&vv)
	case NegInt:
		expect.Negint = (*expectNegint)(&vv)
	case Slice:
		expect.ByteString = expectByteString(vv)
	case String:
		expect.String = (*expectString)(&vv)
	case List:
		expect.List = expectList{}
		for _, vvv := range vv {
			expect.List = append(expect.List, toExpect(vvv))
		}
	case Map:
		expect.Map = expectMap{}
		for k, vvv := range vv {
			expect.Map[k] = toExpect(vvv)
		}
	case *Tag:
		expect.Tag = &expectTag{
			ID:    vv.ID,
			Value: toExpect(vv.Value),
		}
	case Bool:
		expect.Bool = (*expectBool)(&vv)
	case *Nil:
		expect.Null = &expectNull{}
	case *Undefined:
		expect.Undefined = &expectUndefined{}
	case Float32:
		expect.Float32 = (*expectFloat32)(&vv)
	case Float64:
		expect.Float64 = (*expectFloat64)(&vv)
	default:
		panic("what")
	}
	return expect
}

type expectUint Uint
type expectNegint NegInt
type expectByteString Slice
type expectString String
type expectList []successTestExpect
type expectMap map[string]successTestExpect
type expectTag struct {
	ID    uint64            `json:"id"`
	Value successTestExpect `json:"value"`
}
type expectBool Bool
type expectNull Nil
type expectUndefined Undefined
type expectFloat32 Float32
type expectFloat64 Float64

func (v expectUint) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", v)), nil
}
func (v expectNegint) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", v-1)), nil
}
func (v expectByteString) MarshalJSON() ([]byte, error) {
	p := "["
	for i, vv := range v {
		p += fmt.Sprintf("%d", vv)
		if i < len(v)-1 {
			p += ","
		}
	}
	return []byte(p + "]"), nil
}
func (v expectNull) MarshalJSON() ([]byte, error) {
	return []byte("{}"), nil
}
func (v expectUndefined) MarshalJSON() ([]byte, error) {
	return []byte("{}"), nil
}
func (v expectFloat32) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", math.Float32bits(float32(v)))), nil
}
func (v expectFloat64) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", math.Float64bits(float64(v)))), nil
}
