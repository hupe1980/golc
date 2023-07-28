package sqldb

import (
	"regexp"
	"strings"
	"unicode"
)

// CleanQuery cleans sql query from double white space, comments and leading/trailing spaces.
func CleanQuery(query string) string {
	// remove comments
	query = regexp.MustCompile(`/\*(.*)\*/|\-\-(.*)`).ReplaceAllString(query, "")

	// remove double white space
	query = regexp.MustCompile(`\s+`).ReplaceAllString(query, " ")

	// remove leading/trailing space
	query = strings.TrimSpace(query)

	return query
}

type Parser struct {
	query   string
	lowered string
}

func NewParser(query string) *Parser {
	query = CleanQuery(query)

	return &Parser{
		query:   query,
		lowered: strings.ToLower(query),
	}
}

func (p *Parser) IsSelect() bool {
	return p.lowered[:strings.IndexRune(p.lowered, ' ')] == "select"
}

func (p *Parser) TableNames() []string {
	firstSyntax := p.lowered[:strings.IndexRune(p.lowered, ' ')]

	names := make([]string, 0)

	switch firstSyntax {
	case "update":
		i := strings.Index(p.lowered, strings.ToLower("update")) + len("update") + 1
		return append(names, cleanName(p.after(i)))
	case "insert":
		index := regexp.MustCompile("insert(.*?)into").FindStringIndex(p.lowered)
		return append(names, cleanName(p.after(index[1])))
	case "delete":
		index := regexp.MustCompile("delete(.*?)from").FindStringIndex(p.lowered)
		return append(names, cleanName(p.after(index[1])))
	}

	names = append(names, p.tableNamesByFROM()...)

	indices := regexp.MustCompile(strings.ToLower("join")).FindAllStringIndex(p.lowered, -1)
	for _, index := range indices {
		names = append(names, p.after(index[1]))
	}

	return names
}

func (p *Parser) tableNamesByFROM() []string {
	indices := regexp.MustCompile("from(.*?)(left|inner|right|outer|full)|from(.*?)join|from(.*?)where|from(.*?);|from(.*?)$").FindAllStringIndex(p.lowered, -1)

	names := make([]string, 0)

	for _, index := range indices {
		fromStmt := p.lowered[index[0]:index[1]]
		lastSyntax := fromStmt[strings.LastIndex(fromStmt, " ")+1:]

		var tableStmt string
		if lastSyntax == "from" || lastSyntax == "where" || lastSyntax == "left" ||
			lastSyntax == "right" || lastSyntax == "join" || lastSyntax == "inner" ||
			lastSyntax == "outer" || lastSyntax == "full" {
			tableStmt = p.query[index[0]+len("from")+1 : index[1]-len(lastSyntax)-1]
		} else {
			tableStmt = p.query[index[0]+len("from")+1:]
		}

		for _, name := range strings.Split(tableStmt, ",") {
			names = append(names, cleanName(name))
		}
	}

	return names
}

func (p *Parser) after(iWord int) (atAfter string) {
	iAfter := 0

	for i := iWord; i < len(p.lowered); i++ {
		r := rune(p.lowered[i])
		if unicode.IsLetter(r) && iAfter <= 0 {
			iAfter = i
		}

		if (unicode.IsSpace(r) || unicode.IsPunct(r)) && iAfter > 0 {
			atAfter = p.query[iAfter:i]
			break
		}
	}

	if atAfter == "" {
		atAfter = p.query[iAfter:]
	}

	return
}

func cleanName(name string) string {
	name = strings.Fields(name)[0]
	name = strings.TrimSpace(name)
	name = strings.Trim(name, "`")

	lastRune := name[len(name)-1]
	if lastRune == ';' {
		name = name[:len(name)-1]
	}

	return name
}
