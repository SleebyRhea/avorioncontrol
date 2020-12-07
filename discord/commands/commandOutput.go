package commands

import (
	"avorioncontrol/logger"
	"unicode/utf8"
)

// Page is a single pages worth of content that is less than 1500 characters in
// size to confirm with the Discord embed character limit (2000)
type Page struct {
	Content string
	Index   int
}

// CommandOutput is an object describing a commands output
type CommandOutput struct {
	Title       string
	Description string
	Header      string
	Quoted      bool
	Monospace   bool
	Status      int
	Color       int

	current int
	pages   []*Page
	curr    *Page
	next    *Page
	last    *Page

	mutable bool
	lines   []string

	uuid     string
	loglevel int
}

// Loglevel - Return the current loglevel for the CommandRegistrant
func (o *CommandOutput) Loglevel() int {
	return o.loglevel
}

// SetLoglevel - Set the current loglevel
func (o *CommandOutput) SetLoglevel(l int) {
	o.loglevel = l
	logger.LogInfo(o, sprintf("Setting loglevel to %d", l))
}

// UUID - Return a commands name.
func (o *CommandOutput) UUID() string {
	return "CommandOutput:" + o.uuid
}

// AddLine adds a line of output to the CommandOutput
func (o *CommandOutput) AddLine(line string) {
	if o.mutable {
		logger.LogDebug(o, "Adding line")
		o.lines = append(o.lines, line)
	}
}

// ThisPage returns the current page
func (o *CommandOutput) ThisPage() *Page {
	return o.curr
}

// NextPage returns a string of the next page, or an empty string if either the
// doesn't exist or is empty
func (o *CommandOutput) NextPage() *Page {
	page := o.next
	o.last = o.curr

	max := len(o.pages) - 1
	if o.current < max {
		o.current++
	}

	o.curr = o.pages[o.current]

	if o.current == max {
		o.next = nil
	} else {
		o.next = o.pages[o.current+1]
	}

	logger.LogDebug(o, sprintf("Returning page %d", o.current))
	return page
}

// PreviousPage returns a string of the previuos page, or an empty string if
// there is no previous page, or its empty
func (o *CommandOutput) PreviousPage() *Page {
	page := o.last
	o.next = o.curr

	if o.current > 0 {
		o.current--
	}

	o.curr = o.pages[o.current]

	if o.current == 0 {
		o.last = nil
	} else {
		o.last = o.pages[o.current]
	}

	logger.LogDebug(o, sprintf("Returning page %d", o.current))
	return page
}

// Construct processes the lines in a CommandOutput object and forms a linked
// list of Pages to paginate through on embeds
func (o *CommandOutput) Construct() {
	var prefix string
	o.pages = make([]*Page, 0)
	o.pages = append(o.pages, &Page{Content: "", Index: 0})

	i := 0

	if o.Monospace {
		logger.LogDebug(o, "Using monospace output")
	}

	if o.Quoted {
		logger.LogDebug(o, "Using quoted output")
		prefix += "> "
	}

	for lineIndex, line := range o.lines {
		if o.pages[i].Content == "" && o.Monospace {
			o.pages[i].Content = prefix + "```\n"
		}

		if utf8.RuneCountInString(o.pages[i].Content+line) < 1000 {
			o.pages[i].Content += prefix + line + "\n"
		} else {
			if o.Monospace {
				o.pages[i].Content += prefix + "```"
			}

			i++
			o.pages = append(o.pages, &Page{Content: prefix + line + "\n", Index: i})
		}

		logger.LogDebug(o, sprintf("Added line %d to page %d: "+line, lineIndex, i))
	}

	logger.LogDebug(o, sprintf("Processed %d lines", len(o.lines)))
	logger.LogDebug(o, sprintf("Processed %d pages", len(o.pages)))

	if o.Monospace {
		o.pages[len(o.pages)-1].Content += prefix + "```"
	}

	o.mutable = false
	o.current = 0
	o.last = nil
	o.next = nil
	o.curr = o.pages[o.current]

	if len(o.pages) > 1 {
		o.next = o.pages[1]
	}
}

// RemoveLine removes the last line of output from the CommandOutput object
func (o *CommandOutput) RemoveLine() {
	if o.mutable {
		o.pages = o.pages[:len(o.pages)]
	}
}

// Index returns the current index, and the max page index
func (o *CommandOutput) Index() (int, int) {
	return o.current, len(o.pages) - 1
}

func newCommandOutput(cmd *CommandRegistrant, title string) *CommandOutput {
	return &CommandOutput{
		lines:       make([]string, 0),
		Title:       title,
		Description: cmd.description,
		uuid:        cmd.UUID(),
		mutable:     true,
		loglevel:    cmd.Loglevel()}
}
