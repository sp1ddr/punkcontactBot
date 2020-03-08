package punkbot

import (
	"fmt"
)

// MarkdownTags : list of tags that need escape
var MarkdownTags map[rune]string
var foundCloseTag bool
var escapedText []byte

func init() {
	MarkdownTags = make(map[rune]string, 5)

	MarkdownTags['_'] = "\\_"
	MarkdownTags['*'] = "\\*"
	MarkdownTags['`'] = "\\`"

}

// EscapeMarkdown : Escape markdown from string
func EscapeMarkdown(text string) (string, error) {
	if text == "" || len(text) <= 1 {
		return text, nil
	}
	escapedText = make([]byte, len(text))
	foundCloseTag = false
	copy(escapedText, text)

	for pos, char := range text {
		if MarkdownTags[char] != "" {

			if foundCloseTag {
				continue
			}
			fmt.Println("Found md tag => " + string(char))
			// actualPos := pos
			actualTag := char
			fmt.Println("Rest of text ==> " + text[pos+1:len(text)])
			for _, closingTag := range text[pos+1 : len(text)] {

				if len(text[pos+1:len(text)]) == 0 {
					foundCloseTag = false
					break
				}

				if actualTag == closingTag {
					fmt.Println("Found close tag")
					foundCloseTag = true
					break
				}
				foundCloseTag = false
			}

			// If closing tag was not found, escape initial tag
			if !foundCloseTag {
				// Change tag to escaped tag
				fmt.Println("Found unclosed tag => " + string(text[pos]))
				escapedValue := MarkdownTags[char]
				fmt.Println("EscapedValue ==> " + string([]byte(escapedValue)))
				newChar := []byte(escapedValue)
				// have to move every char one position forward
				// or append new char to bit position
				escapedText = append(escapedText[:pos], newChar...)
				escapedText = append(escapedText, []byte(text[pos+1:])...)
			}
			foundCloseTag = false
		}
	}
	fmt.Println(string(escapedText))
	return string(escapedText), nil
}
