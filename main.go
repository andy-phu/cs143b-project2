package main

import (
	"bufio"
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"strconv"
	"strings"
)

type Command struct {
	Type string
	Args []string
}

type disk [][]int

type SegmentTable struct {
	Size        int
	FrameNumber int
}

var log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Caller().Logger()
var initCounter = 0
var prevInit = 0

//=============================================================================================================

func runCommand(c Command, PM *[]int, DISK *disk, ST *[]SegmentTable) {
	cmd := c.Type
	input := c.Args

	switch cmd {
	case "TA":
		if input[0] != cmd {
			log.Info().Msgf("the input of TA is: %v\n", input)
			pa := translate(PM, DISK, ST, input[0])
			log.Info().Msgf("translated: %v", pa)
			fmt.Printf("%v ", pa)
		}
	case "RP":
		if input[0] != cmd {
			log.Info().Msgf("the input of RP is: %v\n", input)
			word := read(PM, input[0])
			fmt.Printf("%v ", word)

		}
	case "NL":
		if input[0] == cmd {
			log.Info().Msgf("the input of NL is: %v\n", input[0])
			fmt.Printf("\n")
		}
	case "PrintInit":
		log.Info().Msgf("PM: %v", *PM) // Print the entire PM array
	default:
		log.Info().Msgf("this is the init file")
		log.Info().Msgf("input1: %s | input2: %s | input3: %s", cmd, input[0], input[1]) // Print the entire PM array
		initial(PM, DISK, ST, cmd, input[0], input[1])
		initCounter++
	}

}

func initial(PM *[]int, DISK *disk, ST *[]SegmentTable, a string, b string, c string) {
	s, _ := strconv.Atoi(a)
	z, _ := strconv.Atoi(b)
	f, _ := strconv.Atoi(c)
	log.Info().Msgf("InitCounter: %d", initCounter)
	//first line of init file
	if initCounter == 0 {
		(*PM)[2*s] = z
		(*PM)[(2*s)+1] = f
		prevInit = f
	} else { //second line
		temp := prevInit * 512
		inner := temp + z
		log.Info().Msgf("Prev Init: %d | Inner: %d", prevInit, inner)
		(*PM)[inner] = f
	}

}

func read(PM *[]int, input string) int {
	pa, _ := strconv.Atoi(input)
	if pa > 524288 {
		return -1
	}
	return (*PM)[pa]
}
func translate(PM *[]int, DISK *disk, ST *[]SegmentTable, input string) int {
	input = strings.TrimSpace(input)

	va, err := strconv.Atoi(input)

	pMask := uint32(0x1FF)
	wMask := uint32(0x1FF)
	pwMask := uint32(0x3FFFF)

	if err != nil { // Handle potential errors
		log.Error().Msgf("Error converting string to int:", err)
		return -1
	}

	s := uint32(va) >> 18
	p := (uint32(va) >> 9) & pMask
	w := uint32(va) & wMask
	pw := uint32(va) & pwMask

	log.Info().Msgf("input: %s | va: %d | s: %d | p: %d | w: %d | pw: %d", input, va, s, p, w, pw)

	size := (*PM)[2*int(s)]

	if int(pw) >= size {
		log.Error().Msgf("PW > PM[2%d]: ", s)
		return -1
	}

	left := ((*PM)[(2*int(s))+1] * 512) + int(p)
	right := int(w)
	element := (*PM)[left]
	quotient := (element * 512) + right
	log.Info().Msgf("left: %d | right: %d | element: %d | quotient: %d", left, right, element, quotient)

	return quotient
}

func main() {
	file, err := os.OpenFile(
		"myapp.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	log = zerolog.New(file).With().Timestamp().Logger()

	scanner := bufio.NewScanner(os.Stdin)
	PM := make([]int, 524288)
	DISK := make(disk, 1024)
	ST := make([]SegmentTable, 1024)

	for i := range PM {
		PM[i] = 0 // Explicitly assign zero to each element
	}

	for i := range ST {
		ST[i] = SegmentTable{Size: -1, FrameNumber: -1}
	}

	for i := range DISK {
		DISK[i] = make([]int, 512)
	}

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Split(line, " ")
		value := parts
		if len(parts) < 2 {
			log.Info().Msgf("No input either NL or init file")
		} else {
			value = parts[1:]
		}
		command := Command{Type: parts[0], Args: value}

		runCommand(command, &PM, &DISK, &ST)
	}

}
