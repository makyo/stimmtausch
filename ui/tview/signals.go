package tview

import (
	"github.com/makyo/stimmtausch/help"
)

// listen listens for events from the signal environment, then does nothing (but
// does it splendidly)
func (t *ui) listen() {
	for {
		res := <-t.listener
		switch res.Name {
		case ">":
			t.views.cycle(1)
			t.update()
		case "<":
			t.views.cycle(-1)
			t.update()
		case "fg":
			// Switch to active view or error if not found.
			if err := t.views.fg(res.Payload[1]); err != nil {
				log.Warningf("error switching connections: %v", err)
			} else {
				t.update()
			}
		case "connect", "c":
			if len(res.Payload) != 1 {
				if t.views.currView != nil {
					if !t.views.currView.connected {
						res.Payload = []string{t.currView.name}
						go t.client.Env.DirectDispatch(res)
					} else {
						log.Warningf("already connected!")
					}
				} else {
					log.Warningf("no worlds available, must provide connection information")
				}
				continue
			}
			v := view.New(res.Payload[1], t.client.Config.Client.UI.Scrollback)
			t.views.add(v)
			fmt.Fprintf("~ Connecting to %s...", res.Payload[1])
		case "disconnect", "dc":
			// If it's a disconnect without a payload, redispatch with the
			// current connection's name.
			if len(res.Payload) != 0 {
				continue
			}
			res.Payload = []string{t.currView.connName}
			log.Tracef("disconnecting current world %+v", res)
			go t.client.Env.DirectDispatch(res)
		case "help":
			// get the command text and tell the system to display it in a modal
			cmd := res.Payload[0]
			h, ok := help.HelpMessages[cmd]
			if !ok {
				log.Warningf("no help available for /%s", cmd)
				continue
			}
			helpText := help.RenderText(h)
			go t.client.Env.Dispatch("_client:showModal", fmt.Sprintf("Help: %s::\n%s", cmd, helpText))
		case "_client:connect":
			// Attach receivedView to conn.
			err := t.connect(res.Payload[0])
			if err != nil {
				log.Errorf("error setting up connection in ui: %v", err)
			}
		case "_client:disconnect":
			// Grey out tab in send title, grey out text in receivedView.
			t.update()
		case "_client:allDisconnect":
			// do we really need to do anything?
			t.updateSendTitle()
		case "_client:showModal":
			res.Name = "_tui:showModal"
			go t.client.Env.DirectDispatch(res)
		case "_tview:showModal":
			fmt.Fprintf(t.currView.buffer, "~~~~~ %s\n%s\n~~~~~\n", res.Payload[0], res.Payload[1])
		default:
			log.Tracef("got unknown signal result %v", res)
			continue
		}
	}
}
