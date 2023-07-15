package helper

import (
	"testing"
)

func TestNormalizePageLimit_ClearInput(t *testing.T) {
	count := int64(100)
	param := map[string]interface{}{
		"offset": int64(10),
		"limit":  10,
	}

	wantPage := int64(2)
	wantLimit := 10
	wantOffset := int64(10)
	page, limit := NormalizePageLimit(count, &param)
	if wantPage != page || wantLimit != limit || wantOffset != param["offset"].(int64) {
		t.Fatalf(`TestNormalizePageLimit_ClearInput = page %d, limit %d, offset %d want match for page %d, limit %d, offset %d`, page, limit, param["offset"], wantPage, wantLimit, wantOffset)
	}
}

func TestNormalizePageLimit_ParamUnset(t *testing.T) {
	count := int64(100)
	param := map[string]interface{}{}

	wantPage := int64(1)
	wantLimit := -1
	var wantOffset int64
	page, limit := NormalizePageLimit(count, &param)
	_, offsetExist := param["offset"]
	if wantPage != page || wantLimit != limit || offsetExist {
		t.Fatalf(`TestNormalizePageLimit_ParamUnset = page %d, limit %d, offset %d want match for page %d, limit %d, offset %d`, page, limit, param["offset"], wantPage, wantLimit, wantOffset)
	}
}

func TestNormalizePageLimit_PageOver(t *testing.T) {
	count := int64(100)
	param := map[string]interface{}{
		"offset": int64(120),
		"limit":  10,
	}

	wantPage := int64(11)
	wantLimit := 10
	wantOffset := int64(100)
	page, limit := NormalizePageLimit(count, &param)
	if wantPage != page || wantLimit != limit || wantOffset != param["offset"].(int64) {
		t.Fatalf(`TestNormalizePageLimit_PageOver = page %d, limit %d, offset %d want match for page %d, limit %d, offset %d`, page, limit, param["offset"], wantPage, wantLimit, wantOffset)
	}
}

func TestNormalizePageLimit_ZeroCount(t *testing.T) {
	count := int64(0)
	param := map[string]interface{}{
		"offset": int64(0),
		"limit":  10,
	}

	wantPage := int64(1)
	wantLimit := 10
	wantOffset := int64(0)
	page, limit := NormalizePageLimit(count, &param)
	if wantPage != page || wantLimit != limit || wantOffset != param["offset"].(int64) {
		t.Fatalf(`TestNormalizePageLimit_ZeroCount = page %d, limit %d, offset %d want match for page %d, limit %d, offset %d`, page, limit, param["offset"], wantPage, wantLimit, wantOffset)
	}
}
