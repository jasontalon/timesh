package log

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	_ "github.com/lib/pq"
)

type CommitLog struct {
	CommitID   string
	Branch     string
	Subject    string
	Company    string
	Project    string
	CommitedAt string
}

func Migrate() error {
	db, err := sql.Open("postgres", os.Getenv("TIMESH_PG_CONN_URI"))

	if err != nil {
		return err
	}
	row := db.QueryRow("select exists (select * from information_schema.tables where table_name = 'commit_logs')")

	exists := false
	if err = row.Scan(&exists); err != nil {
		return err
	}

	if exists {
		return nil
	}

	createTable := `
		CREATE TABLE commit_logs (
			CommitID VARCHAR(150) PRIMARY KEY,
			Company VARCHAR (100),
			Project VARCHAR (100),
			Branch VARCHAR (150),
			Subject VARCHAR (255),
			CommitedAt TIMESTAMPTZ,
			CreatedAt TIMESTAMPTZ
		);`

	_, err = db.Exec(createTable)

	if err != nil {
		return err
	}

	return nil
}

func BulkInsert(c *[]CommitLog, company, project string) (err error) {
	query := `set time zone UTC;
	;with data(CommitID, Company, Project, Branch, Subject, CommitedAt)  as (values %data%) 
	 insert into commit_logs (CommitID, Company, Project, Branch, Subject, CommitedAt, CreatedAt) 
	 select d.CommitID, d.Company, d.Project, d.Branch, d.Subject, d.CommitedAt::timestamptz, timezone('utc', now()) from data d where not exists (select 1 from commit_logs c2
					   where c2.CommitID = d.CommitID);`

	inserts := ""

	for _, item := range *c {
		subject := strings.ReplaceAll(strings.ReplaceAll(item.Subject, `'`, ""), `"`, "")
		inserts += fmt.Sprintf(`('%s','%s','%s','%s','%s','%s'),`, item.CommitID, company, project, item.Branch, subject, item.CommitedAt)
	}

	inserts = strings.TrimRight(inserts, ",")
	query = strings.Replace(query, "%data%", inserts, 1)

	db, err := sql.Open("postgres", os.Getenv("TIMESH_PG_CONN_URI"))

	if err != nil {
		return err
	}

	_, err = db.Exec(query)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func GetCommitLogs(dir, email, since string) (commits []CommitLog, err error) {
	cmd := exec.Command("git", "--no-pager", "log", "--all", "--date", "iso", "--no-merges", "--author", email, `--pretty=format:%ad||%H||%S||%s\n`, "--since", since)
	if dir != "" {
		cmd.Dir = dir
	}

	out, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), `\n`)

	commits = []CommitLog{}

	for _, line := range lines {
		line := strings.TrimRight(line, "\n")
		arr := strings.Split(line, "||")

		if len(arr) != 4 {
			continue
		}

		commit := CommitLog{}
		commit.CommitedAt = arr[0]
		commit.CommitID = arr[1]
		commit.Branch = arr[2]
		commit.Subject = arr[3]

		commits = append(commits, commit)

	}
	return commits, nil
}

func GetAndSave(dir, company, project, email, since string) (err error) {
	commits, err := GetCommitLogs(dir, email, since)

	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("%d commits found", len(commits)))

	if len(commits) == 0 {
		return nil
	}

	err = BulkInsert(&commits, company, project)

	if err != nil {
		return err
	}

	return nil
}
