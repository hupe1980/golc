package textsplitter

type CharacterTextSplitterOptions struct {
	Options
	Separator string
}

type CharacterTextSplitter struct {
	*BaseTextSplitter
	opts CharacterTextSplitterOptions
}

func NewCharacterTextSplitter(optFns ...func(o *CharacterTextSplitterOptions)) *CharacterTextSplitter {
	opts := CharacterTextSplitterOptions{
		Separator: "\n\n",
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

	ts := &CharacterTextSplitter{
		opts: opts,
	}

	ts.BaseTextSplitter = NewBaseTextSplitter(ts.splitText, func(o *Options) {
		o.ChunkSize = opts.ChunkSize
		o.ChunkOverlap = opts.ChunkOverlap
		o.KeepSeparator = opts.KeepSeparator
	})

	return ts
}

func (ts *CharacterTextSplitter) splitText(text string) []string {
	splits := splitTextWithRegex(text, ts.opts.Separator, ts.opts.KeepSeparator)

	separator := ts.opts.Separator
	if ts.opts.KeepSeparator {
		separator = ""
	}

	return ts.mergeSplits(splits, separator)
}
