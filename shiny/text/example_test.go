// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text_test

import (
	"fmt"
	"image"
	"os"

	"golang.org/x/exp/shiny/text"
	"golang.org/x/image/math/fixed"
)

// toyFace implements the font.Face interface by measuring every rune's width
// as 1 pixel.
type toyFace struct{}

func (toyFace) Close() error {
	return nil
}

func (toyFace) Glyph(dot fixed.Point26_6, r rune) (image.Rectangle, image.Image, image.Point, fixed.Int26_6, bool) {
	panic("unimplemented")
}

func (toyFace) GlyphBounds(r rune) (fixed.Rectangle26_6, fixed.Int26_6, bool) {
	panic("unimplemented")
}

func (toyFace) GlyphAdvance(r rune) (fixed.Int26_6, bool) {
	return fixed.I(1), true
}

func (toyFace) Kern(r0, r1 rune) fixed.Int26_6 {
	return 0
}

func printFrame(f *text.Frame, softReturnsOnly bool) {
	for p := f.FirstParagraph(); p != nil; p = p.Next(f) {
		for l := p.FirstLine(f); l != nil; l = l.Next(f) {
			for b := l.FirstBox(f); b != nil; b = b.Next(f) {
				s := b.Text(f)
				if softReturnsOnly {
					if len(s) > 0 && s[len(s)-1] == '\n' {
						s = s[:len(s)-1]
					}
				}
				os.Stdout.Write(s)
			}
			if softReturnsOnly {
				fmt.Println()
			}
		}
	}
}

func Example() {
	var f text.Frame
	f.SetFace(toyFace{})
	// TODO: honor SetMaxWidth, i.e. implement re-layout.
	f.SetMaxWidth(fixed.I(60))

	c := f.NewCaret()
	c.WriteString(mobyDick)
	c.Close()

	fmt.Println("====")
	printFrame(&f, false)
	fmt.Println("====")
	printFrame(&f, true)
	fmt.Println("====")

	// Output:
	// ====
	// CHAPTER 1. Loomings.
	// Call me Ishmael. Some years ago—never mind how long precisely—having little or no money in my purse, and nothing particular to interest me on shore, I thought I would sail about a little and see the watery part of the world...
	// ====
	// CHAPTER 1. Loomings.
	// Call me Ishmael. Some years ago—never mind how long precisely—having little or no money in my purse, and nothing particular to interest me on shore, I thought I would sail about a little and see the watery part of the world...
	//
	// ====
}

const mobyDick = "CHAPTER 1. Loomings.\nCall me Ishmael. Some years ago—never mind how long precisely—having little or no money in my purse, and nothing particular to interest me on shore, I thought I would sail about a little and see the watery part of the world...\n"