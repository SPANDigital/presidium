package validate

type LinkListener = func(link Link)

type Link struct {
	Uri        string
	Location   string
	Status     Status
	Message    string
	IsExternal bool
	Label      string
}

type Report struct {
	Data       map[Status][]Link
	Valid      int // How many valid links we have
	Broken     int // How many broken links we have
	External   int // How many external links we have
	Warning    int // How many warning links we have
	TotalLinks int // The total number of links processed
}

type Status string

const (
	Valid    = Status("valid")
	Broken   = Status("broken")
	Warning  = Status("warning")
	External = Status("external")
)
