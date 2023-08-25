package sqltime_test

import (
	"errors"
	"testing"

	"github.com/utilyre/sqltime"
)

func TestParse(t *testing.T) {
	time := &sqltime.Time{}
	if err := time.Parse("08:56:00"); err != nil {
		t.Error(err)
	}
	if time.Hour != 8 {
		t.Errorf("time.Hour = %d; want 8", time.Hour)
	}
	if time.Minute != 56 {
		t.Errorf("time.Minute = %d; want 56", time.Minute)
	}
	if time.Second != 0 {
		t.Errorf("time.Second = %d; want 0", time.Second)
	}

	if err := time.Parse("18:05"); err != nil {
		t.Error(err)
	}
	if time.Hour != 18 {
		t.Errorf("time.Hour = %d; want 18", time.Hour)
	}
	if time.Minute != 5 {
		t.Errorf("time.Minute = %d; want 5", time.Minute)
	}
	if time.Second != 0 {
		t.Errorf("time.Second = %d; want 0", time.Second)
	}

	if err := time.Parse("20"); err != nil {
		t.Error(err)
	}
	if time.Hour != 20 {
		t.Errorf("time.Hour = %d; want 20", time.Hour)
	}
	if time.Minute != 0 {
		t.Errorf("time.Minute = %d; want 0", time.Minute)
	}
	if time.Second != 0 {
		t.Errorf("time.Second = %d; want 0", time.Second)
	}

	err := time.Parse("24:18:59")
	if err == nil {
		t.Error("time.Parse did not error")
	}
	if !errors.Is(err, sqltime.ErrRange) {
		t.Error("time.Parse did not error sqltime.ErrRange")
	}
}
