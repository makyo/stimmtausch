package tview

import (
	"fmt"
	"regexp"
	//"strings"

	//"github.com/atotto/clipboard"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var (
	// Temporary solution, so not every function has to handle the selection
	// character placement.
	multiSelectionCharWithSelectionToLeftPattern = regexp.MustCompile(selectionChar + "*" + regexp.QuoteMeta(selRegion) + selectionChar + "*" + regexp.QuoteMeta(endRegion))
)

const (
	selectionChar  = string('\u205F')
	leftRegion     = `["left"]`
	rightRegion    = `["right"]`
	selRegion      = `["selection"]`
	endRegion      = `[""]`
	overflowRegion = `["overflow"]`
	emptyText      = selRegion + selectionChar + endRegion
)

// editor is a simple component that wraps tview.TextView in order to give the
// user minimal text edit functionality.
// We should also figure out a way to show the actual cursor, not just fake one.
type editor struct {
	view *tview.TextView
}

// NewEditor Instanciates a ready to use text editor.
func NewEditor() *editor {
	view := tview.NewTextView()

	view.SetWrap(true)
	view.SetWordWrap(true)
	view.SetRegions(true)
	view.SetScrollable(true)
	view.SetText(emptyText)
	view.Highlight("selection")

	e := &editor{
		view: view,
	}

	e.view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		left := []rune(e.view.GetRegionText("left"))
		right := []rune(e.view.GetRegionText("right"))
		selection := []rune(e.view.GetRegionText("selection"))

		//stub
		left = left
		right = right
		selection = selection

		var newText string
		newText = newText
		switch event.Key() {
		case tcell.KeyLeft, tcell.KeyRight:
			e.handleLeftRight(event)
			return nil
		case tcell.KeyUp, tcell.KeyDown:
			e.handleUpDown(event)
			return nil
		case tcell.KeyCtrlV:
			e.handlePaste()
			return nil
		case tcell.KeyCtrlC:
			e.handleCopy()
			return nil
		case tcell.KeyCtrlX:
			e.handleCut()
			return nil
		case tcell.KeyDelete, tcell.KeyBackspace:
			e.handleBackspaceDelete(event)
			return nil
		}
		return event
	})

	/*
		if event.Key() == tcell.KeyLeft &&
			(event.Modifiers() == tcell.ModShift || event.Modifiers() == tcell.ModNone) {
			expandSelection := (event.Modifiers() & tcell.ModShift) == tcell.ModShift
			if len(left) > 0 {
				newText = leftRegion + string(left[:len(left)-1]) + selRegion

				currentSelection := string(selection)
				if currentSelection == selectionChar {
					currentSelection = ""
				}

				if expandSelection {
					newText = newText + string(left[len(left)-1]) + currentSelection + rightRegion + string(right)
				} else {
					newText = newText + string(left[len(left)-1]) + rightRegion + currentSelection + string(right)
				}

				newText = newText + endRegion
				editor.SetText(newText)
			} else if len(selection) > 0 && !expandSelection {
				if len(right) > 0 {
					newText = selRegion + string(selection[0]) + rightRegion + string(selection[1:]) + string(right) + endRegion
				} else {
					newText = selRegion + string(selection[0]) + rightRegion + string(selection[1:]) + endRegion
				}
				editor.setAndFixText(newText)
			}
		} else if event.Key() == tcell.KeyRight &&
			(event.Modifiers() == tcell.ModShift || event.Modifiers() == tcell.ModNone) {
			newText = leftRegion + string(left)
			expandSelection := (event.Modifiers() & tcell.ModShift) == tcell.ModShift
			if len(right) > 0 {
				if expandSelection {
					newText = newText + selRegion + string(selection) + string(right[0]) + rightRegion + string(right[1:])
				} else {
					newText = newText + string(selection) + selRegion + string(right[0]) + rightRegion + string(right[1:])
				}
			} else {
				endsWithSelectionChar := strings.HasSuffix(string(selection), selectionChar)
				if !endsWithSelectionChar {
					if expandSelection {
						newText = newText + selRegion + string(selection)
					} else if !expandSelection {
						newText = newText + string(selection) + selRegion
					}

					newText = newText + selectionChar
				} else {
					if expandSelection {
						newText = newText + selRegion + string(selection)
					} else {
						newText = newText + string(selection[:len(selection)-1]) + selRegion + selectionChar
					}
				}
			}

			newText = newText + endRegion
			editor.setAndFixText(newText)
		} else if event.Key() == tcell.KeyLeft &&
			(event.Modifiers()&(tcell.ModShift|tcell.ModCtrl)) == (tcell.ModShift|tcell.ModCtrl) {
			if len(left) > 0 {
				selectionFrom := 0
				for i := len(left) - 2; /*Skip space left to selection*/ /* i >= 0; i-- {
				if left[i] == ' ' || left[i] == '\n' {
					selectionFrom = i
					break
				}
			}

			if selectionFrom != 0 {
				newText = leftRegion + string(left[:selectionFrom+1]) + selRegion + string(left[selectionFrom+1:]) + string(string(selection)) + rightRegion + string(right) + endRegion
			} else {
				newText = selRegion + string(left) + string(string(selection)) + rightRegion + string(right) + endRegion
			}
			editor.setAndFixText(newText)
		}
	} else if event.Key() == tcell.KeyRight &&
		(event.Modifiers()&(tcell.ModShift|tcell.ModCtrl)) == (tcell.ModShift|tcell.ModCtrl) {
		if len(right) > 0 {
			selectionFrom := len(right) - 1
			for i := 1; /*Skip space right to selection*/ /* i < len(right)-1; i++ {
				if right[i] == ' ' || right[i] == '\n' {
					selectionFrom = i
					break
				}
			}

			if selectionFrom != len(right)-1 {
				newText = leftRegion + string(left) + selRegion + string(string(selection)) + string(right[:selectionFrom]) + rightRegion + string(right[selectionFrom:]) + endRegion
			} else {
				newText = leftRegion + string(left) + selRegion + string(string(selection)) + string(right) + endRegion
			}
			editor.setAndFixText(newText)
		}
	} else if event.Key() == tcell.KeyRight &&
		event.Modifiers() == tcell.ModCtrl {
		if len(right) > 0 {
			selectionAt := len(right) - 1
			for i := 1; /*Skip space right to selection*/ /* i < len(right)-1; i++ {
				if right[i] == ' ' || right[i] == '\n' {
					selectionAt = i
					break
				}
			}

			if selectionAt != len(right)-1 {
				newText = leftRegion + string(left) + string(string(selection)) + string(right[:selectionAt]) + selRegion + string(right[selectionAt]) + rightRegion + string(right[selectionAt+1:]) + endRegion
			} else {
				newText = leftRegion + string(left) + string(selection) + string(right) + selRegion + selectionChar + endRegion
			}
			editor.setAndFixText(newText)
		}
	} else if event.Key() == tcell.KeyLeft &&
		event.Modifiers() == tcell.ModCtrl {
		if len(left) > 0 {
			selectionAt := 0
			for i := len(left) - 2; /*Skip space left to selection*/ /* i >= 0; i-- {
					if left[i] == ' ' || left[i] == '\n' {
						selectionAt = i
						break
					}
				}

				if selectionAt != 0 {
					newText = leftRegion + string(left[:selectionAt]) + selRegion + string(left[selectionAt]) + rightRegion + string(left[selectionAt+1:]) + string(string(selection)) + string(right) + endRegion
				} else {
					if len(left) > 1 {
						newText = selRegion + string(left[0]) + rightRegion + string(left[1:]) + string(selection) + string(right) + endRegion
					} else {
						newText = selRegion + string(left[0]) + rightRegion + string(selection) + string(right) + endRegion
					}
				}
				editor.setAndFixText(newText)
			}
		} else if event.Key() == tcell.KeyCtrlA {
			if len(left) > 0 || len(right) > 0 {
				newText = selRegion + string(left) + string(selection) + string(right) + endRegion
				editor.setAndFixText(newText)
			}
		} else if event.Key() == tcell.KeyBackspace2 ||
			event.Key() == tcell.KeyBackspace {
			if len(selection) == 1 && len(left) >= 1 {
				newText = leftRegion + string(left[:len(left)-1]) + selRegion + string(selection) + rightRegion + string(right) + endRegion
				editor.SetText(newText)
			} else if len(selection) > 1 {
				newText = leftRegion + string(left) + selRegion
				if len(right) > 0 {
					newText = newText + string(right[0]) + rightRegion + string(right[1:])
				} else {
					newText = newText + selectionChar
				}
				newText = newText + endRegion
				editor.setAndFixText(newText)
			}
		} else if event.Key() == tcell.KeyDelete {
			if len(selection) >= 1 && strings.HasSuffix(string(selection), selectionChar) {
				newText = leftRegion + string(left) + selRegion + selectionChar + endRegion
				editor.setAndFixText(newText)
			} else if string(selection) != selectionChar {
				newText = leftRegion + string(left) + selRegion
				if len(right) == 0 {
					newText = newText + selectionChar
				} else {
					newText = newText + string(right[0])
				}

				if len(right) > 1 {
					newText = newText + rightRegion + string(right[1:])
				}

				newText = newText + endRegion
				editor.setAndFixText(newText)
			}
		} else {
			var character rune
			if shortcuts.EventsEqual(event, shortcuts.InputNewLine.Event) {
				character = '\n'
			} else if !shortcuts.EventsEqual(event, shortcuts.SendMessage.Event) {
				character = event.Rune()
			}

			if event.Key() == tcell.KeyCtrlV {
				if editor.inputCapture != nil {
					result := editor.inputCapture(event)
					if result == nil {
						return nil
					}
				}

				clipBoardContent, clipError := clipboard.ReadAll()
				if clipError == nil {
					if string(selection) == selectionChar {
						newText = leftRegion + string(left) + clipBoardContent + selRegion + string(selection)
					} else {
						newText = leftRegion + string(left) + clipBoardContent
						if len(selection) == 1 {
							newText = newText + selRegion + string(selection) + rightRegion + string(right)
						} else {
							newText = newText + selRegion
							if len(right) == 0 {
								newText = newText + selectionChar
							} else if len(right) == 0 {
								newText = newText + string(right[0])
							} else {
								newText = newText + string(right[0]) + rightRegion + string(right[1:])
							}
						}
					}
					editor.setAndFixText(newText + endRegion)
					editor.triggerHeightRequestIfNeccessary()
				}
				return nil
			}

			if character == 0 && editor.inputCapture != nil {
				editor.inputCapture(event)
				return nil
			}

			if len(right) == 0 {
				if len(selection) == 1 {
					if string(selection) == selectionChar {
						editor.setAndFixText(fmt.Sprintf(`["left"]%s%s[""]["selection"]%s[""]`, string(left), (string)(character), string(selectionChar)))
					} else {
						editor.setAndFixText(fmt.Sprintf(`["left"]%s%s[""]["selection"]%s[""]`, string(left), (string)(character), string(selection)))
					}
				} else {
					editor.setAndFixText(fmt.Sprintf(`["left"]%s%s[""]["selection"]%s[""]`, string(left), (string)(character), string(selectionChar)))
				}
			} else {
				editor.setAndFixText(fmt.Sprintf(`["left"]%s%s[""]["selection"]%s[""]["right"]%s[""]`,
					string(left), string(character), string(selection), string(right)))
			}
		}

		atIndex := -1
		newLeft := editor.GetRegionText("left")
		for i := len(newLeft) - 1; i >= 0; i-- {
			if newLeft[i] == ' ' {
				break
			}

			if newLeft[i] == '@' {
				atIndex = i
				break
			}
		}

		editor.ScrollToHighlight()

		return nil
	}*/

	return e
}

func (e *editor) GetPrimitive() tview.Primitive {
	return e.view
}

func (e *editor) handleLeftRight(event *tcell.EventKey) {
}

func (e *editor) handleUpDown(event *tcell.EventKey) {
}

func (e *editor) handlePaste() {
}

func (e *editor) handleCopy() {
}

func (e *editor) handleCut() {
}

func (e *editor) handleBackspaceDelete(event *tcell.EventKey) {
}

func (editor *editor) setAndFixText(text string) {
	newText := multiSelectionCharWithSelectionToLeftPattern.ReplaceAllString(text, selRegion+selectionChar+endRegion)
	editor.setViewText(newText)
}

// SetText sets the texts of the base TextView, but also sets the selection
// and necessary groups for the navigation behaviour.
func (e *editor) setViewText(text string) {
	if text == "" {
		e.view.SetText(emptyText)
	} else {
		e.view.SetText(fmt.Sprintf("[\"left\"]%s[\"\"][\"selection\"]%s[\"\"]", text, string(selectionChar)))
	}
}

// getViewText returns the text without color tags, region tags and so on.
func (e *editor) getViewText() string {
	left := e.view.GetRegionText("left")
	right := e.view.GetRegionText("right")
	selection := e.view.GetRegionText("selection")

	if right == "" && selection == string(selectionChar) {
		return left
	}

	return left + selection + right
}
