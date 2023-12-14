package prompt

import (
	"bytes"
	"regexp"
	"text/template"
	"text/template/parse"
)

type FormatterOptions struct {
	IgnoreMissingKeys bool
	TemplateFuncMap   template.FuncMap
}

type Formatter struct {
	text     string
	template *template.Template
	fields   []string
}

func NewFormatter(text string, optFns ...func(o *FormatterOptions)) *Formatter {
	opts := FormatterOptions{
		IgnoreMissingKeys: false,
		TemplateFuncMap:   make(map[string]any),
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	t := template.Must(template.New("template").Funcs(opts.TemplateFuncMap).Parse(text))

	if !opts.IgnoreMissingKeys {
		t = t.Option("missingkey=error")
	}

	return &Formatter{
		text:     text,
		template: t,
		fields:   ListTemplateFields(t),
	}
}

func (pt *Formatter) Render(values map[string]any) (string, error) {
	var doc bytes.Buffer
	if err := pt.template.Execute(&doc, values); err != nil {
		return "", err
	}

	return doc.String(), nil
}

func (pt *Formatter) Fields() []string {
	return pt.fields
}

func ListTemplateFields(t *template.Template) []string {
	return listNodeFields(t.Tree.Root)
}

func listNodeFields(node parse.Node) []string {
	res := []string{}
	if node.Type() == parse.NodeAction {
		res = append(res, node.String())
	}

	if ln, ok := node.(*parse.ListNode); ok {
		for _, n := range ln.Nodes {
			res = append(res, listNodeFields(n)...)
		}
	}

	return res
}

func extractNameFromField(input string) string {
	re := regexp.MustCompile(`{{\.(.*?)}}`)
	matches := re.FindStringSubmatch(input)

	if len(matches) == 2 {
		return matches[1]
	}

	return ""
}
