package main

import (
	"reflect"
	"testing"
)

func TestParseMainActionWithThreeField(t *testing.T) {
	// 	var ParseMainActionTests = []struct {
	// 		in       string // input
	// 		expected Action // expected result
	// }{
	// 		{"<p>Missing Attendance record 欠缺出入勤紀錄:<br /> 2023-12-21<br />2023-12-28<br />2024-01-08<br />2024-01-11<br />2024-01-12<br />2024-01-15</p><p>Early leave:<br /> 2023-12-18<br />2024-01-17</p><p>Lateness 遲到:<br /> 2023-12-18<br />2024-01-04</p>", 1},
	// }

	threeFieldString := "<p>Missing Attendance record 欠缺出入勤紀錄:<br /> 2023-12-21<br />2023-12-28<br />2024-01-08<br />2024-01-11<br />2024-01-12<br />2024-01-15</p><p>Early leave:<br /> 2023-12-18<br />2024-01-17</p><p>Lateness 遲到:<br /> 2023-12-18<br />2024-01-04</p>"

	want := &Action{
		missAttendance: make([]string, 0, 31),
		earlyLeave:     make([]string, 0, 31),
		lateness:       make([]string, 0, 31),
	}
	want.missAttendance = append(want.missAttendance, "2023-12-21", "2023-12-28", "2024-01-08", "2024-01-11", "2024-01-12", "2024-01-15")
	want.earlyLeave = append(want.earlyLeave, "2023-12-18", "2024-01-17")
	want.lateness = append(want.lateness, "2023-12-18", "2024-01-04")

	got := ParseMainAction(threeFieldString)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParseMainAction(threeFieldString) = %s, want %s", got, want)
	}
}
