package pinecone

import (
	"fmt"
)

type Endpoint struct {
	IndexName   string
	ProjectName string
	Environment string
}

func (e *Endpoint) String() string {
	return fmt.Sprintf("%s-%s.svc.%s.pinecone.io:443", e.IndexName, e.ProjectName, e.Environment)
}
