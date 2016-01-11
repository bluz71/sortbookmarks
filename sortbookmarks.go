package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	re "regexp"
	"sort"
	st "strings"
)

var (
	bookmarkFormat string
	bookmarks      = make([]folder, 0, 32)
	bookmarkHeader = `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<!-- This is an automatically generated file.
     It will be read and overwritten.
     DO NOT EDIT! -->
<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>
<DL><p>
    <DT><H3 ADD_DATE="1431142608" LAST_MODIFIED="1431142616" PERSONAL_TOOLBAR_FOLDER="true">Bookmarks Bar</H3>
    <DL><p>`
	bookmarkFooter = `    </DL><p>
</DL><p>`
	subfolderRe = re.MustCompile(`\s{12}<DT><H3 ADD_DATE="\d*" LAST_MODIFIED="\d*">(\w+)<\/H3>`)
	folderRe    = re.MustCompile(`\s{8}<DT><H3 ADD_DATE="\d*" LAST_MODIFIED="\d*">(\w+)<\/H3>`)
	pageRe      = re.MustCompile(`\s+<DT><A HREF="([\w:?&=\-\/\.]+)" ADD_DATE="\d+".*>(.+)<\/A>`)
)

type page struct {
	name string
	url  string
}

type folder struct {
	name       string
	subfolders []folder
	pages      []page
}

type pageByName []page
type folderByName []folder

func (slice pageByName) Len() int { return len(slice) }
func (slice pageByName) Less(i, j int) bool {
	return st.ToLower(slice[i].name) < st.ToLower(slice[j].name)
}
func (slice pageByName) Swap(i, j int) { slice[i], slice[j] = slice[j], slice[i] }

func (slice folderByName) Len() int { return len(slice) }
func (slice folderByName) Less(i, j int) bool {
	return st.ToLower(slice[i].name) < st.ToLower(slice[j].name)
}
func (slice folderByName) Swap(i, j int) { slice[i], slice[j] = slice[j], slice[i] }

func readAndProcessBookmarks() {
	file, err := os.Open("bookmarks.html")
	if err != nil {
		log.Fatal("'bookmarks.html' no longer exists.")
	}
	defer file.Close()

	var f, sf *folder
	addToSubfolder := false

	// Go through bookmarks file line by line.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		matches := subfolderRe.FindStringSubmatch(scanner.Text())
		if matches != nil {
			f.subfolders = append(
				f.subfolders,
				folder{matches[1], make([]folder, 0), make([]page, 0)})
			sf = &f.subfolders[len(f.subfolders)-1]
			addToSubfolder = true
			continue
		}
		matches = folderRe.FindStringSubmatch(scanner.Text())
		if matches != nil {
			bookmarks = append(
				bookmarks,
				folder{matches[1], make([]folder, 0), make([]page, 0)})
			f = &bookmarks[len(bookmarks)-1]
			addToSubfolder = false
			continue
		}
		matches = pageRe.FindStringSubmatch(scanner.Text())
		if matches != nil {
			if addToSubfolder {
				sf.pages = append(sf.pages, page{matches[2], matches[1]})
			} else {
				f.pages = append(f.pages, page{matches[2], matches[1]})
			}
			continue
		}
	}
}

func printBookmarks(bms []folder, indent string) {
	sort.Sort(folderByName(bms))
	for _, f := range bms {
		fmt.Println(indent + "        <DT><H3 ADD_DATE=\"1414814291\" LAST_MODIFIED=\"1432192423\">" +
			f.name + "</H3>")
		fmt.Println(indent + "        <DL><p>")
		printBookmarks(f.subfolders, "    ")
		sort.Sort(pageByName(f.pages))
		for _, p := range f.pages {
			fmt.Println(indent + "            <DT><A HREF=\"" + p.url +
				"\" ADD_DATE=\"1251002296\">" + p.name + "</A>")
		}
		fmt.Println(indent + "        </DL><p>")
	}
}

func printMozillaBookmarks() {
	fmt.Println(bookmarkHeader)
	printBookmarks(bookmarks, "")
	fmt.Println(bookmarkFooter)
}

func main() {
	// sortbookmarks must be called with an argument.
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "Usage: sortbookmarks mozilla")
		os.Exit(2)
	}

	// Only the 'mozilla' argument is currently accepted.
	if os.Args[1] != "mozilla" {
		fmt.Fprintln(os.Stderr, "Unknown sortbookmarks argument:",
			os.Args[1], ", expecting 'mozilla'.")
		os.Exit(2)
	}
	bookmarkFormat = os.Args[1]

	_, err := os.Stat("bookmarks.html")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Expecting 'bookmarks.html' in current directory.")
		os.Exit(2)
	}

	readAndProcessBookmarks()
	printMozillaBookmarks()
}
