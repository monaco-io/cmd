package font

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"path"
	"strconv"
	"strings"
)

const (
	defaultFont = "epic"
	defaultText = "ascii-art"
)

//go:embed fonts/*
var fonts embed.FS

func List() {
	for _, v := range getAllFaceNames() {
		fmt.Println(v)
	}
}

func getAllFaceNames() []string {
	files, _ := fonts.ReadDir("fonts")
	var names []string
	for i := range files {
		names = append(names, strings.Replace(files[i].Name(), ".flf", "", 1))
	}
	return names
}

func EchoAll(text string) {
	if text == "" {
		text = defaultText
	}
	for _, v := range getAllFaceNames() {
		fmt.Printf("\nFont: %s\n", v)
		Echo(text, v)
	}
}

func fontBytes(name string) []byte {
	b, _ := fonts.ReadFile(path.Join("fonts", name+".flf"))
	return b
}

func AsciiArt(text, fontName string) (str string) {
	if fontName == "" {
		fontName = defaultFont
	}
	font := newFont(fontName)
	var rows []string
	for r := 0; r < font.height; r++ {
		printRow := ""
		for _, char := range text {
			fontIndex := char - 32
			printRow += font.letters[fontIndex][r]
		}
		if len(strings.TrimSpace(printRow)) > 0 {
			rows = append(rows, strings.TrimRight(printRow, " "))
		}
	}
	for _, printRow := range rows {
		str += fmt.Sprintf("%s\n", printRow)
	}
	return
}

func Echo(text, fontName string) {
	fmt.Println(AsciiArt(text, fontName))
}

type font struct {
	height  int
	letters [][]string
}

func newFont(name string) (font font) {
	fontBytesReader := bytes.NewReader(fontBytes(name))
	scanner := bufio.NewScanner(fontBytesReader)
	font.setAttributes(scanner)
	font.setLetters(scanner)
	return font
}

func (font *font) setAttributes(scanner *bufio.Scanner) {
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "flf2") {
			datum := strings.Fields(text)[1]
			font.height, _ = strconv.Atoi(datum)
			break
		}
	}
}

func (font *font) setLetters(scanner *bufio.Scanner) {
	font.letters = append(font.letters, make([]string, font.height))
	for i := range font.letters[0] {
		font.letters[0][i] = "  "
	}
	letterIndex := 0
	for scanner.Scan() {
		text, cutLength, letterIndexInc := scanner.Text(), 1, 0
		if lastCharLine(text, font.height) {
			font.letters = append(font.letters, []string{})
			letterIndexInc = 1
			if font.height > 1 {
				cutLength = 2
			}
		}
		if letterIndex > 0 {
			appendText := ""
			if len(text) > 1 {
				appendText = text[:len(text)-cutLength]
			}
			font.letters[letterIndex] = append(font.letters[letterIndex], appendText)
		}
		letterIndex += letterIndexInc
	}
}

func lastCharLine(text string, height int) bool {
	endOfLine, length := "  ", 2
	if height == 1 && len(text) > 0 {
		length = 1
	}
	if len(text) >= length {
		endOfLine = text[len(text)-length:]
	}
	return endOfLine == strings.Repeat("@", length) ||
		endOfLine == strings.Repeat("#", length) ||
		endOfLine == strings.Repeat("$", length)
}
