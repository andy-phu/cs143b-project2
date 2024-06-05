package main

import (
	"bufio"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/zyedidia/generic/list"
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
var freeList = list.New[int]()

// =============================================================================================================
func removeValue(value int) {
	log.Info().Msgf("FREE LIST REMOVAL: %d", value)
	cur := freeList.Front
	for cur != nil {
		if cur.Value == value {
			freeList.Remove(cur)
		}
		cur = cur.Next
	}
}

func runCommand(c Command, PM *[]int, DISK *disk, ST *[]SegmentTable) bool {
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
	case "PrintDisk":
		log.Info().Msgf("DISK: %v", *DISK) // Print the entire PM array
	default:
		log.Info().Msgf("this is the init file")
		//log.Info().Msgf("input1: %s | input2: %s | input3: %s", cmd, input[0], input[1]) // Print the entire PM array
		log.Info().Msgf("cmd: %s | input: %v", cmd, input)
		combined := make([]string, 0, len(input)+1)
		combined = append(combined, cmd)
		combined = append(combined, input...)

		if !initial(PM, DISK, ST, combined) {
			log.Error().Msgf("ERROR: size of input array is not a multiple of 3 (%d)", len(input))
			return false
		}

		removeValue(initCounter)
		initCounter++
	}
	return true

}

func initial(PM *[]int, DISK *disk, ST *[]SegmentTable, input []string) bool {
	//Line 1: 8 4000 3 9 5000 ‐7
	//every 3 is szf, which defines where the segment tables reside in the PM

	//input array has to be multiples of 3
	if len(input)%3 != 0 {
		log.Error().Msgf("the size of the input array is: %d", len(input))
		return false
	}

	for i := 0; i < len(input); i += 3 {

		s, _ := strconv.Atoi(input[i])
		z, _ := strconv.Atoi(input[i+1])
		f, _ := strconv.Atoi(input[i+2])

		log.Info().Msgf("InitCounter: %d", initCounter)
		//first line of init file
		if initCounter == 0 {
			(*PM)[2*s] = z
			(*PM)[(2*s)+1] = f
			removeValue(f)
		} else { //second line
			prevInit := (*PM)[(2*s)+1]

			//page fault because the
			if prevInit < 0 {
				(*DISK)[-1*prevInit][z] = f
				log.Info().Msgf("page fault: %v", prevInit)
			} else {
				temp := prevInit * 512
				inner := temp + z
				log.Info().Msgf("Prev Init: %d | Inner: %d", prevInit, inner)
				(*PM)[inner] = f
				//log.Info().Msgf("Not a Page fault and the frame is no longer free : %v", f)
				if f > 0 {
					removeValue(f)
				}
			}
		}
	}

	return true
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

	if err != nil || va > 524288 { // Handle potential errors
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

	ptIndex := (*PM)[(2*int(s))+1]
	log.Info().Msgf("ptIndex: %v", ptIndex)
	finalFrame := 0
	if ptIndex < 0 {
		prevPtIndex := ptIndex
		head := freeList.Front
		freeFrameIndex := head.Value
		log.Info().Msgf("freeFrameIndex: %v", freeFrameIndex)
		removeValue(freeFrameIndex)
		//allocate the new free frame to the ptIndex
		//and make the physical memory have the page table at newpt * 512
		ptIndex = freeFrameIndex
		block := (*DISK)[-1*prevPtIndex][int(p)]
		finalFrame = block
		log.Info().Msgf("putting new frame index here : %v", (2*int(s))+1)
		(*PM)[(2*int(s))+1] = ptIndex
		log.Info().Msgf("BLOCK : %v at x: %d | y: %d", block, -1*prevPtIndex, int(p))
		//transfer over prevPtIndex from disk to the new one
		//(*PM)[ptIndex*512+int(p)] = block
		for i, blockVal := range (*DISK)[-1*prevPtIndex] {
			(*PM)[ptIndex*512+int(p)+i] = blockVal
		}
	}

	pgIndex := (*PM)[(ptIndex*512)+int(p)]
	log.Info().Msgf("pgIndex: %v", pgIndex)
	if pgIndex < 0 {
		head := freeList.Front
		freeFrameIndex := head.Value
		finalFrame = freeFrameIndex
		removeValue(freeFrameIndex)
		(*PM)[(ptIndex*512)+int(p)] = freeFrameIndex
	}

	if finalFrame == 0 {
		return pgIndex*512 + int(w)
	}
	quotient := (finalFrame * 512) + int(w)

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

	log = zerolog.New(file).With().Logger()

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

	for i := 0; i < 1024; i++ {
		freeList.PushBack(i)
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

		if !runCommand(command, &PM, &DISK, &ST) {
			break
			log.Error().Msgf("Error executing command")
		}
	}

}
