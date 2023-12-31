package sqltime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	time := Time{}

	if err := time.Parse("08:56:04"); assert.Nil(t, err) {
		assert.Equal(t, Time{hour: 8, minute: 56, second: 4}, time)
	}
	if err := time.Parse("9:10"); assert.Nil(t, err) {
		assert.Equal(t, Time{hour: 9, minute: 10, second: 0}, time)
	}
	if err := time.Parse("06"); assert.Nil(t, err) {
		assert.Equal(t, Time{hour: 6, minute: 0, second: 0}, time)
	}

	if err := time.Parse("7:30:13:4"); assert.NotNil(t, err) {
		assert.ErrorIs(t, err, ErrTooManyParts)
	}

	if err := time.Parse("hi:5:47"); assert.NotNil(t, err) {
		target := new(AtoiError)
		if assert.ErrorAs(t, err, &target) {
			assert.Equal(t, target.Part, PartHour)
		}
	}
	if err := time.Parse("5:hi:47"); assert.NotNil(t, err) {
		target := new(AtoiError)
		if assert.ErrorAs(t, err, &target) {
			assert.Equal(t, target.Part, PartMinute)
		}
	}
	if err := time.Parse("5:47:hi"); assert.NotNil(t, err) {
		target := new(AtoiError)
		if assert.ErrorAs(t, err, &target) {
			assert.Equal(t, target.Part, PartSecond)
		}
	}

	if err := time.Parse("24:57:00"); assert.NotNil(t, err) {
		target := new(RangeError)
		if assert.ErrorAs(t, err, &target) {
			assert.Equal(t, RangeError{Part: PartHour, Min: 0, Max: 23}, *target)
		}
	}
	if err := time.Parse("13:-2:06"); assert.NotNil(t, err) {
		target := new(RangeError)
		if assert.ErrorAs(t, err, &target) {
			assert.Equal(t, RangeError{Part: PartMinute, Min: 0, Max: 59}, *target)
		}
	}
	if err := time.Parse("20:33:62"); assert.NotNil(t, err) {
		target := new(RangeError)
		if assert.ErrorAs(t, err, &target) {
			assert.Equal(t, RangeError{Part: PartSecond, Min: 0, Max: 59}, *target)
		}
	}
}

func TestAdd(t *testing.T) {
	assert.Equal(
		t,
		Time{hour: 15, minute: 18, second: 45},
		Time{hour: 13, minute: 10, second: 45}.Add(Time{hour: 2, minute: 8, second: 0}),
	)

	assert.Equal(
		t,
		Time{hour: 0, minute: 1, second: 30},
		Time{hour: 23, minute: 5, second: 20}.Add(Time{hour: 0, minute: 56, second: 10}),
	)

	assert.Equal(
		t,
		Time{hour: 15, minute: 0, second: 22},
		Time{hour: 14, minute: 59, second: 32}.Add(Time{hour: 0, minute: 0, second: 50}),
	)
}

func TestSub(t *testing.T) {
	assert.Equal(
		t,
		Time{hour: 11, minute: 2, second: 45},
		Time{hour: 13, minute: 10, second: 45}.Sub(Time{hour: 2, minute: 8, second: 0}),
	)

	assert.Equal(
		t,
		Time{hour: 22, minute: 59, second: 59},
		Time{hour: 23, minute: 4, second: 20}.Sub(Time{hour: 0, minute: 4, second: 21}),
	)

	assert.Equal(
		t,
		Time{hour: 23, minute: 59, second: 28},
		Time{hour: 0, minute: 0, second: 2}.Sub(Time{hour: 0, minute: 0, second: 34}),
	)
}
