package dataframe

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDataframeBasic(t *testing.T) {
	expectScriptOutput(t, "testdata/dataframe_basic.star", "testdata/dataframe_basic.expect.txt")
}

func TestDataframeBoolSelect(t *testing.T) {
	expectScriptOutput(t, "testdata/dataframe_bool_select.star",
		"testdata/dataframe_bool_select.expect.txt")
}

func TestDataframeBoolSelectDontUseEqualOperator(t *testing.T) {
	_, err := runScript(t, "testdata/dataframe_bool_select_failure.star")
	if err == nil {
		t.Fatal("error expected, did not get one")
	}
	expectErr := "cannot call DataFrame.Get with bool. If you are trying `df[df[column] == val], instead use `df[df[column].equals(val)]`"
	if err.Error() != expectErr {
		t.Errorf("error mismatch\nwant: %s\ngot: %s", expectErr, err)
	}
}

func TestDataframeSize(t *testing.T) {
	expectScriptOutput(t, "testdata/dataframe_size.star", "testdata/dataframe_size.expect.txt")
}

func TestDataframeConcat(t *testing.T) {
	expectScriptOutput(t, "testdata/dataframe_concat.star", "testdata/dataframe_concat.expect.txt")
}

func TestDataframeAt(t *testing.T) {
	expectScriptOutput(t, "testdata/dataframe_at.star", "testdata/dataframe_at.expect.txt")
}

func TestDataframeSetKey(t *testing.T) {
	expectScriptOutput(t, "testdata/dataframe_setkey.star", "testdata/dataframe_setkey.expect.txt")
}

func TestDataframeGetKey(t *testing.T) {
	expectScriptOutput(t, "testdata/dataframe_getkey.star", "testdata/dataframe_getkey.expect.txt")
}

func TestDataframeApply(t *testing.T) {
	expectScriptOutput(t, "testdata/dataframe_apply.star", "testdata/dataframe_apply.expect.txt")
}

func TestDataframeAppend(t *testing.T) {
	expectScriptOutput(t, "testdata/dataframe_append.star", "testdata/dataframe_append.expect.txt")
}

func TestDataframeDropDuplicates(t *testing.T) {
	expectScriptOutput(t, "testdata/dataframe_drop_duplicates.star",
		"testdata/dataframe_drop_duplicates.expect.txt")
}

func TestDataframeHead(t *testing.T) {
	expectScriptOutput(t, "testdata/dataframe_head.star", "testdata/dataframe_head.expect.txt")
}

func TestDataframeMerge(t *testing.T) {
	expectScriptOutput(t, "testdata/dataframe_merge.star", "testdata/dataframe_merge.expect.txt")
}

func TestDataframeReadCSV(t *testing.T) {
	expectScriptOutput(t, "testdata/dataframe_read_csv.star",
		"testdata/dataframe_read_csv.expect.txt")
}

func TestDataframeResetIndex(t *testing.T) {
	expectScriptOutput(t, "testdata/dataframe_reset_index.star",
		"testdata/dataframe_reset_index.expect.txt")
}

func TestDataframeColumns(t *testing.T) {
	expectScriptOutput(t, "testdata/dataframe_columns.star", "testdata/dataframe_columns.expect.txt")
}

func TestDataframeColumnsNone(t *testing.T) {
	runTestScript(t, "testdata/dataframe_columns_none.star",
		"testdata/dataframe_columns_none.expect.txt")
}

func TestDataframeGroupBy(t *testing.T) {
	expectScriptOutput(t, "testdata/dataframe_groupby.star", "testdata/dataframe_groupby.expect.txt")
}

func TestDataframeNotImplemented(t *testing.T) {
	_, err := runScript(t, "testdata/dataframe_not_implemented.star")
	if err == nil {
		t.Fatal("error expected, did not get one")
	}
	expectErr := `dataframe.ffill is not implemented. If you need this functionality to exist, file an issue at 'https://github.com/qri-io/starlib/issues' with the title 'dataframe.ffill needs implementation'. Please first search if an issue exists already`
	if err.Error() != expectErr {
		t.Errorf("error mismatch\nwant: %s\ngot: %s", expectErr, err)
	}
}

type invalidData struct{}

func TestDataframeFromRows(t *testing.T) {
	// Construct a valid dataframe from a single row of various types of data
	rows := [][]interface{}{}
	record := []interface{}{"test", 31.2, 11.4, "ok", int64(597), "", 107, 6.91}
	rows = append(rows, record)
	df, err := NewDataFrame(rows, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := df.String()
	expect := `        0     1     2   3    4  5    6    7
0    test  31.2  11.4  ok  597    107  6.9
`
	if got != expect {
		t.Errorf("mismatch: expect %q, got %q", expect, got)
	}

	// Ensure a dataframe with an invalid type ends up returning an error
	rows = [][]interface{}{}
	record = []interface{}{"test", 31.2, &invalidData{}}
	rows = append(rows, record)
	_, err = NewDataFrame(rows, nil, nil)
	if err == nil {
		t.Fatal("expected to get an error, did not get one")
	}
	expectErr := `invalid object &{} of type *dataframe.invalidData`
	if expectErr != err.Error() {
		t.Errorf("error mismatch, expect: %s, got: %s", expectErr, err)
	}

	// Construct a dataframe from multiple rows
	rows = [][]interface{}{}
	record = []interface{}{"test", 31.2, 17, int64(45)}
	rows = append(rows, record)
	record = []interface{}{"more", 9.8, 62, int64(3)}
	rows = append(rows, record)
	df, err = NewDataFrame(rows, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	got = df.String()
	expect = `        0     1   2   3
0    test  31.2  17  45
1    more   9.8  62   3
`
	if got != expect {
		t.Errorf("mismatch: expect %q, got %q", expect, got)
	}

	// Construct a dataframe with non-matching columns, they get casted correctly
	rows = [][]interface{}{}
	record = []interface{}{"test", 31.2, 17, int64(45)}
	rows = append(rows, record)
	record = []interface{}{25, "ok", int64(4), "hi"}
	rows = append(rows, record)
	df, err = NewDataFrame(rows, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	got = df.String()
	expect = `        0     1   2   3
0    test  31.2  17  45
1      25    ok   4  hi
`
	if got != expect {
		t.Errorf("mismatch: expect %q, got %q", expect, got)
	}
}

func TestDataframeFromList(t *testing.T) {
	ls := []interface{}{1.2, 3.4, 5.6}
	df, err := NewDataFrame(ls, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := df.String()
	expect := `       0
0    1.2
1    3.4
2    5.6
`
	if got != expect {
		t.Errorf("mismatch: expect %q, got %q", expect, got)
	}
}

func TestDataframeFromSeries(t *testing.T) {
	s := newSeriesFromObjects([]interface{}{"a", "b", "c"}, nil, "")
	df, err := NewDataFrame(s, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := df.String()

	expect := `     0
0    a
1    b
2    c
`
	if got != expect {
		t.Errorf("mismatch: expect %q, got %q", expect, got)
	}
}

type someStruct struct {
	ID     int
	Name   string
	Sounds []string
}

func TestDataframeAccessor(t *testing.T) {
	// Construct a dataframe with a few rows and columns
	rows := [][]interface{}{
		[]interface{}{"test", 31.2, 11.4, "ok", int64(597), "", 107, 6.91},
		[]interface{}{"more", 7.8, 44.1, "hi", int64(612), "", 94, 3.1},
		[]interface{}{"last", 90.2, 26.8, "yo", int64(493), "", 272, 4.3},
	}
	df, err := NewDataFrame(rows, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Print it and test the result
	got := df.String()
	expectText := `        0     1     2   3    4  5    6    7
0    test  31.2  11.4  ok  597    107  6.9
1    more   7.8  44.1  hi  612     94  3.1
2    last  90.2  26.8  yo  493    272  4.3
`
	if diff := cmp.Diff(expectText, got); diff != "" {
		t.Errorf("mismatch (-want +got):%s\n", diff)
	}

	// Retrieve a single element
	actual, err := df.At2d(1, 6)
	if err != nil {
		t.Fatal(err)
	}

	// Check that it is correct
	expectNum := 94
	if diff := cmp.Diff(expectNum, actual); diff != "" {
		t.Errorf("mismatch (-want +got):%s\n", diff)
	}

	// Modify an element
	if err := df.SetAt2d(2, 3, "ah"); err != nil {
		t.Fatal(err)
	}

	// Print it and test the result has been modified
	got = df.String()
	expectText = `        0     1     2   3    4  5    6    7
0    test  31.2  11.4  ok  597    107  6.9
1    more   7.8  44.1  hi  612     94  3.1
2    last  90.2  26.8  ah  493    272  4.3
`
	if diff := cmp.Diff(expectText, got); diff != "" {
		t.Errorf("mismatch (-want +got):%s\n", diff)
	}

	// Modify an element by assigning a structured object
	structObj := someStruct{
		ID:     1,
		Name:   "cat",
		Sounds: []string{"meow", "purr"},
	}
	if err := df.SetAt2d(0, 3, structObj); err != nil {
		t.Fatal(err)
	}

	// Retrieve the struct element, and type convert it
	actual, err = df.At2d(0, 3)
	if err != nil {
		t.Fatal(err)
	}
	actualObj, ok := actual.(someStruct)
	if !ok {
		t.Fatalf("expected to retrieve a someStruct{}, got %v", actual)
	}

	// Check that it is correct
	if actualObj.ID != 1 {
		t.Errorf("expectd ID == 1, got %v", actualObj.ID)
	}
	if actualObj.Name != "cat" {
		t.Errorf("expected Name == cat, got %v", actualObj.Name)
	}
	expectSounds := []string{"meow", "purr"}
	if diff := cmp.Diff(expectSounds, actualObj.Sounds); diff != "" {
		t.Errorf("mismatch (-want +got):%s\n", diff)
	}
}
