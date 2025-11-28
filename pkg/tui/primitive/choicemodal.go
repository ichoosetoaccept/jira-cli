package primitive

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ChoiceModal is a modal dialog that presents a list of choices to the user.
type ChoiceModal struct {
	*tview.Box
	frame  *tview.Frame
	text   string
	list   *tview.List
	footer *tview.TextView
	done   func(index int, label string)
}

// Choice modal layout constants.
const (
	choiceFooterHeight = 2
)

// NewChoiceModal creates a new choice modal.
func NewChoiceModal() *ChoiceModal {
	m := &ChoiceModal{Box: tview.NewBox()}

	m.list = tview.NewList().
		ShowSecondaryText(false).
		SetMainTextColor(tcell.ColorDefault)

	m.footer = tview.NewTextView()
	m.footer.SetTitleAlign(tview.AlignCenter)
	m.footer.SetTextAlign(tview.AlignCenter)
	m.footer.SetTextStyle(tcell.StyleDefault.Italic(true))
	m.footer.SetBorderPadding(1, 0, 0, 0)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(m.list, 0, 1, true).
		AddItem(m.footer, choiceFooterHeight, 0, false)

	m.frame = tview.NewFrame(flex).SetBorders(0, 0, 1, 0, 0, 0)
	m.frame.SetBorder(true).SetBorderPadding(1, 1, 1, 1)

	return m
}

// SetText sets the text displayed in the modal.
func (m *ChoiceModal) SetText(text string) {
	m.text = text
}

// SetDoneFunc sets the callback function when a choice is selected.
func (m *ChoiceModal) SetDoneFunc(doneFunc func(index int, label string)) *ChoiceModal {
	m.done = doneFunc
	return m
}

// SetChoices sets the list of choices to display.
func (m *ChoiceModal) SetChoices(choices []string) *ChoiceModal {
	m.list.Clear()
	for _, choice := range choices {
		m.list.AddItem(choice, "", 0, nil)
	}
	return m
}

// SetSelected sets the currently selected choice index.
func (m *ChoiceModal) SetSelected(index int) *ChoiceModal {
	m.list.SetCurrentItem(index)
	return m
}

// GetFooter returns the footer text view.
func (m *ChoiceModal) GetFooter() *tview.TextView {
	return m.footer
}

// Focus is called when this primitive receives focus.
func (m *ChoiceModal) Focus(delegate func(p tview.Primitive)) {
	delegate(m.list)
}

// HasFocus returns whether or not this primitive has focus.
func (m *ChoiceModal) HasFocus() bool {
	return m.list.HasFocus()
}

// InputHandler returns the handler for this primitive.
func (m *ChoiceModal) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return m.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		switch event.Key() {
		case tcell.KeyEnter:
			if m.done != nil {
				index := m.list.GetCurrentItem()
				label, _ := m.list.GetItemText(index)
				m.done(index, label)
			}
		default:
			if handler := m.frame.InputHandler(); handler != nil {
				handler(event, setFocus)
			}
		}
	})
}

// MouseHandler returns the mouse handler for this primitive.
func (m *ChoiceModal) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (bool, tview.Primitive) {
	return m.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (bool, tview.Primitive) {
		if handler := m.frame.MouseHandler(); handler != nil {
			return handler(action, event, setFocus)
		}
		return false, nil
	})
}

// Choice modal Draw constants.
const (
	choiceVerticalMargin   = 3
	choiceFrameExtraHeight = 7
	choiceModalWidth       = 70
	marginMultiplier       = 2
)

// Draw draws this primitive onto the screen.
func (m *ChoiceModal) Draw(screen tcell.Screen) {
	screenWidth, screenHeight := screen.Size()
	width := choiceModalWidth

	m.frame.Clear()
	var lines []string
	for _, line := range strings.Split(m.text, "\n") {
		if line == "" {
			lines = append(lines, line)
			continue
		}
		lines = append(lines, tview.WordWrap(line, width)...)
	}

	for _, line := range lines {
		m.frame.AddText(line, true, tview.AlignCenter, tcell.ColorDefault)
	}

	height := len(lines) + m.list.GetItemCount() + choiceFrameExtraHeight
	maxHeight := screenHeight - choiceVerticalMargin*marginMultiplier
	if height > maxHeight {
		height = maxHeight
	}

	x := (screenWidth - width) / centerDivisor
	y := (screenHeight - height) / centerDivisor
	m.frame.SetRect(x, y, width, height)
	m.frame.Draw(screen)
}
