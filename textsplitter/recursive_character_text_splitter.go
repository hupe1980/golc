package textsplitter

import (
	"regexp"
)

type RecursiveCharacterTextSplitterOptions struct {
	Options
	Separators []string
}

type RecursiveCharacterTextSplitter struct {
	*BaseTextSplitter
	opts RecursiveCharacterTextSplitterOptions
}

func NewRecusiveCharacterTextSplitter(optFns ...func(o *RecursiveCharacterTextSplitterOptions)) *RecursiveCharacterTextSplitter {
	opts := RecursiveCharacterTextSplitterOptions{
		Separators: []string{"\n\n", "\n", " ", ""},
		Options: Options{
			ChunkSize:     4000,
			ChunkOverlap:  200,
			KeepSeparator: false,
			LengthFunc: func(text string) int {
				return len(text)
			},
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	ts := &RecursiveCharacterTextSplitter{
		opts: opts,
	}

	ts.BaseTextSplitter = NewBaseTextSplitter(ts.splitText, func(o *Options) {
		o.ChunkSize = opts.ChunkSize
		o.ChunkOverlap = opts.ChunkOverlap
		o.KeepSeparator = opts.KeepSeparator
	})

	return ts
}

func (ts *RecursiveCharacterTextSplitter) splitText(text string) []string {
	return ts.splitTextBySeparators(text, ts.opts.Separators)
}

func (ts *RecursiveCharacterTextSplitter) splitTextBySeparators(text string, separators []string) []string {
	finalChunks := make([]string, 0)
	separator := separators[len(separators)-1]
	newSeparators := make([]string, 0)

	for i, s := range ts.opts.Separators {
		if s == "" {
			separator = s
			break
		}

		if regexp.MustCompile(s).MatchString(text) {
			separator = s
			newSeparators = separators[i+1:]

			break
		}
	}

	splits := splitTextWithRegex(text, separator, ts.opts.KeepSeparator)
	goodSplits := make([]string, 0)
	separatorToMerge := ""

	if !ts.opts.KeepSeparator {
		separatorToMerge = separator
	}

	for _, s := range splits {
		if ts.opts.LengthFunc(s) < ts.opts.ChunkSize {
			goodSplits = append(goodSplits, s)
		} else {
			if len(goodSplits) > 0 {
				mergedText := ts.mergeSplits(goodSplits, separatorToMerge)
				finalChunks = append(finalChunks, mergedText...)
				goodSplits = nil
			}

			if len(newSeparators) == 0 {
				finalChunks = append(finalChunks, s)
			} else {
				otherInfo := ts.splitTextBySeparators(s, newSeparators)
				finalChunks = append(finalChunks, otherInfo...)
			}
		}
	}

	if len(goodSplits) > 0 {
		mergedText := ts.mergeSplits(goodSplits, separatorToMerge)
		finalChunks = append(finalChunks, mergedText...)
	}

	return finalChunks
}
