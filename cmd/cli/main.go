package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	timesh "github.com/jasontalon/timesh/pkg/log"
)

func main() {
	email, since, company, project, dir := flag.String("email", "", "author of email"),
		flag.String("since", "yesterday", "range of commit"),
		flag.String("company", "", "company name"),
		flag.String("project", "", "project name"),
		flag.String("dir", "", "directory")

	flag.Parse()
	if *email == "" {
		fmt.Println("email is required")
		return
	}

	if *company == "" {
		fmt.Println("company is required")
		return
	}

	if *project == "" {
		fmt.Println("project is required")
		return
	}

	if con := os.Getenv("TIMESH_PG_CONN_URI"); con == "" {
		log.Fatal("env TIMESH_PG_CONN_URI does not exists")
		return
	}

	if err := timesh.Migrate(); err != nil {
		log.Fatal(err)
		return
	}

	if err := timesh.GetAndSave(*dir, *company, *project, *email, *since); err != nil {
		log.Fatal(err)
		return
	}
}
