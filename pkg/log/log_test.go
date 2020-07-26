package log

import (
	"fmt"
	"testing"

	_ "github.com/lib/pq"
)

func TestMigrateTable(t *testing.T) {
	if err := Migrate(); err != nil {
		t.Fatal(err)
	}
}

func TestDisplayGitLogs(t *testing.T) {
	commits, err := GetCommitLogs("../../asurion/mgpd-optus-web-video", "jason.talon@asurion.com", "1 weeks ago")

	if err != nil {
		t.Fatal(err)
		return
	}

	fmt.Println(commits)
}

func TestGetAndSave(t *testing.T) {
	err := GetAndSave("../../asurion/mgpd-optus-web-video", "Asurion", "Optus", "jason.talon@asurion.com", "3 days ago")
	if err != nil {
		t.Fatal(err)
	}
}
