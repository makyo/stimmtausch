package charm

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/makyo/stimmtausch/client"
	"github.com/makyo/stimmtausch/signal"
)

type model struct {
	client        *client.Client
	errs          *History
	currViewIndex int
	inputs        textarea.Model
	sent          []*History
	views         []viewport.Model
	listener      chan signal.Signal
	ready         chan bool
}

type errMsg struct {
	err error
}

type inputMsg string

var log = loggo.GetLogger("stimmtausch.ui")

func listen(c *client.Client) tea.Msg {
	return func() tea.Msg{
	}
}

func (m model) Init() tea.Cmd {
	return listen(m.c)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case errMsg:
		log.Errorf("UI encountered an error: %v", errMsg.err)
		return m, listen(m.c)
	case inputMsg:
		log.Errorf("not implemented - input: %s", inputMsg)
		return m, listen(m.c)
	}
	return m, listen(m.c)
}

func (m model) View() (string) {
}

func (m model) Run(c *client.Client, done, ready chan bool) {
	log.Tracef("creating UI")
	m = model{
		client: c,
		sent: []*History{}
		errs: NewHistory(100, true)
		ready: ready
	}
	p := tea.NewProgram(m, tea.WithAltView())
	if _, err := p.Run(); err != nil {
		log.Criticalf("unable to create UI: %v", err)
		os.Exit(2)
	}
	done <- true
}
