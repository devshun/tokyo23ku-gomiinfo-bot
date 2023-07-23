package model

import (
	"testing"
)

func TestGarbageType_String(t *testing.T) {
	tests := []struct {
		name string
		g    GarbageType
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.String(); got != tt.want {
				t.Errorf("GarbageType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWeekday_String(t *testing.T) {
	tests := []struct {
		name string
		w    Weekday
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.w.String(); got != tt.want {
				t.Errorf("Weekday.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
