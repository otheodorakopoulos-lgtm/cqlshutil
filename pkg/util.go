package pkg

const (
	ListCommand     = "list"
	DownloadCommand = "download"
	GreaterThenFlag = "gt"
	LesserThenFlag  = "lt"
	OutputFileFlag  = "o"
)

// Downloadable defines the filter criteria
type Downloadable struct {
	FullName string
	Version  string
}

// Version represents a cqlsh release and its cloud availability
type Version struct {
	Name       string
	CloudState string
}
