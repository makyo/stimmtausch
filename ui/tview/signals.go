package tview

import (
	"fmt"

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
			go t.app.QueueUpdateDraw(t.update)
		case "<":
			t.views.cycle(-1)
			go t.app.QueueUpdateDraw(t.update)
		case "fg":
			// Switch to active view or error if not found.
			if err := t.views.fg(res.Payload[1]); err != nil {
				log.Warningf("error switching connections: %v", err)
			} else {
				go t.app.QueueUpdateDraw(t.update)
			}
		case "connect", "c":
			if len(res.Payload) != 1 {
				if t.views.currView != nil {
					if !t.views.currView.connected {
						res.Payload = []string{t.views.currView.name}
						go t.client.Env.DirectDispatch(res)
					} else {
						log.Warningf("already connected!")
					}
				} else {
					log.Warningf("no worlds available, must provide connection information")
				}
				continue
			}
		case "disconnect", "dc":
			// If it's a disconnect without a payload, redispatch with the
			// current connection's name.
			if len(res.Payload) != 0 {
				continue
			}
			res.Payload = []string{t.views.currView.name}
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
			conn, ok := t.client.Conn(res.Payload[0])
			if !ok {
				log.Errorf("unable to find connection %s", res.Payload[0])
				continue
			}
			v := NewView(t.app, res.Payload[0], t.client.Config.Client.UI.Scrollback)
			t.views.add(v)
			fmt.Fprintf(v, "[-:-:d]~ Connecting to %s...[-:-:-]\n\n", res.Payload[0])
			v.conn = conn
			v.displayName = conn.GetDisplayName()
			conn.AddOutput(fmt.Sprintf("tview-%d", t.views.currIndex), v.buffer, true)
			err := conn.Open()
			if err != nil {
				log.Errorf("error setting up connection in ui: %v", err)
			}
			v.connected = conn.Connected
			go t.app.QueueUpdateDraw(t.update)
		case "_client:disconnect":
			// Grey out tab in send title, grey out text in receivedView.
			go t.app.QueueUpdateDraw(t.update)
		case "_client:allDisconnect":
			// do we really need to do anything?
			go t.app.QueueUpdateDraw(t.update)
		case "_client:showModal":
			res.Name = "_tui:showModal"
			go t.client.Env.DirectDispatch(res)
		case "_tview:showModal":
			fmt.Fprintf(t.views.currView.buffer, "~~~~~ %s\n%s\n~~~~~\n", res.Payload[0], res.Payload[1])
		default:
			log.Tracef("got unknown signal result %v", res)
			continue
		}
	}
}
