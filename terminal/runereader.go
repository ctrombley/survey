package terminal

import (
	"fmt"
	"os"
)

type RuneReader struct {
	Input *os.File

	state runeReaderState
}

func NewRuneReader(input *os.File) *RuneReader {
	return &RuneReader{
		Input: input,
		state: newRuneReaderState(input),
	}
}

func (rr *RuneReader) ReadLine(mask rune) ([]rune, error) {
	line := []rune{}

	// we only care about horizontal displacements from the origin so start counting at 0
	index := 0

	for {
		// wait for some input
		r, _, err := rr.ReadRune()
		if err != nil {
			return line, err
		}

		// if the user pressed enter or some other newline/termination like ctrl+d
		if r == '\r' || r == '\n' || r == KeyEndTransmission {
			// go to the beginning of the next line
			Print("\r\n")

			// we're done processing the input
			return line, nil
		}

		// if the user interrupts (ie with ctrl+c)
		if r == KeyInterrupt {
			// go to the beginning of the next line
			Print("\r\n")

			// we're done processing the input, and treat interrupt like an error
			return line, fmt.Errorf("interrupt")
		}

		// allow for backspace/delete editing of password
		if r == KeyBackspace || r == KeyDelete {
			// and we're not at the beginning of the line
			if index > 0 {

				if len(line) > 0 {
					line = line[:len(line)-1]
				}
				CursorBack(1)
				EraseLine(ERASE_LINE_END)
				index--
			}
			// decrement the index
			continue
		}

		// if the left arrow is pressed
		if r == KeyArrowLeft {
			// and we have space to the left
			if index > 0 {
				// move the cursor to the left
				CursorBack(1)
				// decrement the index
				index--

			} else {
				// otherwise we are at the beginning of where we started reading lines

				// sound the bell?
			}

			// we're done processing this key press
			continue
		}

		// if the right arrow is pressed
		if r == KeyArrowRight {
			// and we have space to the right of the word
			if index < len(line) {
				// move the cursor to the right
				CursorForward(1)
				// increment the index
				index++

			} else {
				// otherwise we are at the end of the word and can't go past

				// sound the bell?
			}

			// we're done processing this key press
			continue
		}

		// ignore up and down arrow so they don't add to line
		if r == KeyArrowUp || r == KeyArrowDown {
			continue
		}

		// the user pressed a regular key

		// if we are at the end of the line
		if index == len(line) {
			// just append the character at the end of the line
			line = append(line, r)
			// increment the location counter
			index++

			// if we don't need to mask the input
			if mask == 0 {
				// just print the character the user pressed
				Printf("%c", r)
			} else {
				// otherwise print the mask we were given
				Printf("%c", mask)
			}
		} else {
			// we are in the middle of the word so we need to insert the character the user pressed
			line = append(line[:index], append([]rune{r}, line[index:]...)...)

			// visually insert the character by deleting the rest of the line
			EraseLine(ERASE_LINE_END)

			// print the rest of the word after
			Printf(string(line[index:]))

			// leave the cursor where the user left it
			CursorBack(len(line) - index - 1)

			// accomodate the new letter in our counter
			index++
		}
	}
}
