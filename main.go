package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-ini/ini"
	"github.com/ysugimoto/cho"
)

// field name constants in .aws/credentials
const (
	accessKey       = "aws_access_key_id"
	secretAccessKey = "aws_secret_access_key"
)

// Credential struct
type credential struct {
	name   string
	key    string
	secret string
}

// format makes string for selectable line
func (c credential) format(offset int) string {
	format := "%-" + fmt.Sprint(offset) + "s: %s | %s"
	return fmt.Sprintf(format, c.name, c.key, c.secret)
}

// toExport makes environment exports string for source import
func (c credential) toExport() io.Reader {
	exports := []string{
		fmt.Sprintf(`export AWS_ACCESS_KEY_ID="%s"`, c.key),
		fmt.Sprintf(`export AWS_SECRET_ACCESS_KEY="%s"`, c.secret),
	}
	return strings.NewReader(strings.Join(exports, "\n") + "\n")
}

// Fixed credential filepath
var credentialFile = filepath.Join(os.Getenv("HOME"), ".aws/credentials")

// main function
func main() {
	if _, err := os.Stat(credentialFile); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	accounts, err := ini.Load(credentialFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Parse and factory to map
	var max int
	creds := map[string]credential{}
	for _, name := range accounts.SectionStrings() {
		if name == "DEFAULT" {
			continue
		}
		if len(name) > max {
			max = len(name)
		}
		section := accounts.Section(name)
		creds[name] = credential{
			name:   name,
			key:    section.Key(accessKey).Value(),
			secret: section.Key(secretAccessKey).Value(),
		}
	}

	var credentialName string
	// If arguments isn't supplied, choose from list
	if len(os.Args) < 2 {
		credentialName = selectFromCredentailList(creds, max)
		// Otherwise, use it
	} else {
		credentialName = os.Args[1]
	}

	v, ok := creds[credentialName]
	if !ok {
		fmt.Printf("Credential %s not found\n", credentialName)
		os.Exit(1)
	}
	io.Copy(os.Stdout, v.toExport())

	// Output debug string
	fmt.Fprintf(os.Stderr, "AWS credential set as %s\n", v.name)
	fmt.Fprintf(os.Stderr, "  AWS_ACCESS_KEY_ID:      %s\n", v.key)
	fmt.Fprintf(os.Stderr, "  AWS_SECRET_ACCESS_KEY:  %s\n", v.secret)
}

// Select credential from list
func selectFromCredentailList(creds map[string]credential, max int) string {
	list := []string{}
	for _, c := range creds {
		list = append(list, c.format(max+1))
	}
	retChan := make(chan string, 1)
	terminate := make(chan struct{})
	go cho.Run(list, retChan, terminate)
	selected := ""
LOOP:
	for {
		select {
		case selected = <-retChan:
			break LOOP
		case <-terminate:
			os.Exit(1)
		}
	}
	spl := strings.SplitN(selected, ":", 2)
	return strings.TrimSpace(spl[0])
}
