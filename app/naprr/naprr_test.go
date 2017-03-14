package main

import (
	"github.com/nsip/nias2/naprr"
	//"log"
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func errcheck(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Error %s", err)
	}
}

func TestFileLoadAndReport(t *testing.T) {

	s, _ := os.Getwd()
	err := os.RemoveAll(filepath.Join(s, "out"))
	if _, err := os.Stat("./out/codeframe.csv"); err == nil {
		t.Fatalf("codeframe.csv not cleared")
	}
	err = os.MkdirAll(filepath.Join(s, "out"), 0777)
	errcheck(t, err)

	filename := "../../unit_test_files/minimaster.xml.zip"
	clearNSSWorkingDirectory()
	di := naprr.NewDataIngest()
	di.RunSynchronous(filename)
	rb := naprr.NewReportBuilder()
	rb.Run()
	rw := naprr.NewReportWriter()
	rw.Run()

	if _, err := os.Stat("./out/codeframe.csv"); os.IsNotExist(err) {
		t.Fatalf("codeframe.csv not generated")
	}
	if _, err := os.Stat("./out/codeframe_writing.csv"); os.IsNotExist(err) {
		t.Fatalf("codeframe_writing.csv not generated")
	}
	if _, err := os.Stat("./out/domain_scores.csv"); os.IsNotExist(err) {
		t.Fatalf("domain_scores.csv not generated")
	}
	if _, err := os.Stat("./out/participation.csv"); os.IsNotExist(err) {
		t.Fatalf("participation.csv not generated")
	}
	if _, err := os.Stat("./out/score_summary.csv"); os.IsNotExist(err) {
		t.Fatalf("score_summary.csv not generated")
	}
	if _, err := os.Stat("./out/21212/domain_scores.csv"); os.IsNotExist(err) {
		t.Fatalf("21212/domain_scores.csv not generated")
	}
}

func TestReportFiles(t *testing.T) {
	filescompare(t, "./out/codeframe.csv", "../../unit_test_files/naprr_out/codeframe.csv")
	filescompare(t, "./out/codeframe_writing.csv", "../../unit_test_files/naprr_out/codeframe_writing.csv")
	filescompare(t, "./out/domain_scores.csv", "../../unit_test_files/naprr_out/domain_scores.csv")
	filescompare(t, "./out/participation.csv", "../../unit_test_files/naprr_out/participation.csv")
	filescompare(t, "./out/score_summary.csv", "../../unit_test_files/naprr_out/score_summary.csv")
	filescompare(t, "./out/score_summary.csv", "../../unit_test_files/naprr_out/score_summary.csv")
	filescompare(t, "./out/21212/domain_scores.csv", "../../unit_test_files/naprr_out/21212/domain_scores.csv")
	filescompare(t, "./out/21212/participation.csv", "../../unit_test_files/naprr_out/21212/participation.csv")
	filescompare(t, "./out/21212/score_summary.csv", "../../unit_test_files/naprr_out/21212/score_summary.csv")
	filescompare(t, "./out/21213/domain_scores.csv", "../../unit_test_files/naprr_out/21213/domain_scores.csv")
	filescompare(t, "./out/21213/participation.csv", "../../unit_test_files/naprr_out/21213/participation.csv")
	filescompare(t, "./out/21213/score_summary.csv", "../../unit_test_files/naprr_out/21213/score_summary.csv")
}

func filescompare(t *testing.T, file1name string, file2name string) {
	file1, err := readLines(file1name)
	errcheck(t, err)
	file2, err := readLines(file2name)
	errcheck(t, err)
	if len(file1) != len(file2) {
		t.Fatalf("%s does not have the expected number of lines\n", file1name)
	}
	sort.Strings(file1)
	sort.Strings(file2)
	for i := 0; i < len(file1); i++ {
		if !strings.EqualFold(file1[i], file2[i]) {
			t.Fatalf("Line in %s is not the expected value:\nExpected: %s\nActual  : %s\n", file1name, file2[i], file1[i])
		}
	}
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func TestMain(m *testing.M) {
	clearNSSWorkingDirectory()
	/* clear out the out directory */
	os.RemoveAll("./out/*/*")
	os.RemoveAll("./out/*.csv")
	ss := launchNatsStreamingServer()
	defer ss.Shutdown()
	os.Exit(m.Run())
}
