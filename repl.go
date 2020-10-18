package main

import (
	"github.com/gdamore/tcell"
)

// TODO: Scroll screen so lookahead cells are always shown to the left and right

// How many ciphertext characters to display after the plaintext
const lookahead = 5

type REPL struct {
	crib   *Crib
	screen tcell.Screen
	cursor Position
}

func NewREPL(crib *Crib) *REPL {
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}

	if err := screen.Init(); err != nil {
		panic(err)
	}

	cursor := Position{0, 0}

	return &REPL{crib, screen, cursor}
}

func (r *REPL) update() bool {
	switch event := r.screen.PollEvent().(type) {
	case *tcell.EventResize:
		// TODO
	case *tcell.EventKey:
		switch event.Key() {
		case tcell.KeyESC:
			r.screen.Fini()
			return false
		case tcell.KeyUp:
			r.moveCursorUp()
		case tcell.KeyDown:
			r.moveCursorDown()
		case tcell.KeyLeft:
			r.moveCursorLeft()
		case tcell.KeyRight:
			r.moveCursorRight()
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			r.moveCursorLeft()
			r.crib.remove(r.cursor)
		default:
			b := byte(event.Rune())
			if r.crib.guess(r.cursor, b) {
				r.moveCursorRight()
			}
		}
	}

	return true
}

func (r *REPL) moveCursorUp() {
	r.cursor.i = max(r.cursor.i - 1, 0)
}

func (r *REPL) moveCursorDown() {
	r.cursor.i = min(r.cursor.i + 1, len(r.crib.ciphertexts) - 1)
}

func (r *REPL) moveCursorLeft() {
	r.cursor.j = max(r.cursor.j - 1, 0)
}

func (r *REPL) moveCursorRight() {
	r.cursor.j = min(r.cursor.j + 1, len(r.crib.key) + lookahead - 1)
}

func (r *REPL) render() {
	r.screen.Clear()

	for i := 0; i < len(r.crib.ciphertexts); i++ {
		for j := 0; j < len(r.crib.key) && j < len(r.crib.ciphertexts[i]); j++ {
			pos := Position{i, j}

			style := tcell.StyleDefault
			if pos == r.cursor {
				style = style.Reverse(true)
			}

			b := r.crib.get(pos)
			r.screen.SetContent(pos.j, pos.i, rune(b), []rune{}, style)
		}

		hex := r.crib.hexCiphertexts[i]
		for j := 0; j < lookahead && 2 * len(r.crib.key) + 2 * j < len(hex); j++ {
			cursorJ := len(r.crib.key) + j

			style := tcell.StyleDefault.Foreground(tcell.ColorRed)
			if i == r.cursor.i && cursorJ == r.cursor.j {
				style = style.Reverse(true)
			}

			screenJ := len(r.crib.key) + 2 * j
			hexJ := 2 * len(r.crib.key) + 2 * j

			r.renderByte(i, screenJ + 0, hex[hexJ + 0], style)
			r.renderByte(i, screenJ + 1, hex[hexJ + 1], style)
		}
	}

	r.screen.Show()
}

func (r *REPL) renderByte(i, j int, b byte, style tcell.Style) {
	r.screen.SetContent(j, i, rune(b), []rune{}, style)
}

