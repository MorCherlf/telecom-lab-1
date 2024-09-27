package main

import(
	"fmt"
)

type BaseStation struct {
	ID int
	Word string
	WalshCode []int
}

func walsh(n int) [][]int {
	if n == 1 {
		return [][]int{{1}}
	}
	half := walsh(n / 2)
	result := make([][]int, n)
	for i := 0; i < n/2; i++ {
		result[i] = append(half[i], half[i]...)
		result[i+n/2] = append(half[i], invert(half[i])...)
	}
	return result
}

func invert(arr []int) []int {
	inv := make([]int, len(arr))
	for i, v := range arr {
		if v == 1 {
			inv[i] = -1
		} else {
			inv[i] = 1
		}
	}
	return inv
}

func wordToASCIIBinary(word string) []int {
	ascii := make([]int, 0)
	for _, char := range word{
		// From letter to ascii to binary as length 8
		for i := 7; i >= 0; i-- {
			bit := (int(char) >> i) & 1
			ascii = append(ascii, bit)
		}
	}
	return ascii
}

// Transfer binary 0 =>> -1
func preprocessing(binary []int) []int {
	for i := 0; i < len(binary); i++ {
		if binary[i] == 0 {
			binary[i] = -1
		}
	}
	return binary
}

func encode(data []int, walshCode []int) []int {
	encoded := []int{}
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(walshCode); j++ {
			encoded = append(encoded, data[i] * walshCode[j])
		}
	}
	return encoded
}

// Use MMSE to detect mixed signal
func mmseDetection(mixedSignals []int, baseStations []BaseStation) []string {
	detectedWords := make([]string, len(baseStations))
	for i, bs := range baseStations {
		estimate := 0
		for j := 0; j < len(mixedSignals); j++ {
			estimate += mixedSignals[j] * mixedSignals[j]
		}
		threshold := 0
		if estimate > threshold {
			detectedWords[i] = bs.Word
		} else {
			detectedWords[i] = ""
		}
	}
	return detectedWords
}

func main() {
	// Generate WalshCode
	walshCodeLength := 8
	code := walsh(walshCodeLength)
	baseStations := []BaseStation{
		{1, "GOD", code[1]},
		{2, "CAT", code[2]},
		{3, "HAM", code[3]},
		{4, "SUN", code[4]},
		{5, "USA", code[5]},
	}

	// Encode Words and transfer to Signals
	timeSteps := 0
	fullSignals := make([][]int, len(baseStations))
	for b,bs := range baseStations {
        ascii := wordToASCIIBinary(bs.Word)
		data := preprocessing(ascii)
		wordLength := len(data)/8
		signal := []int{}
		for i := 0; i < wordLength; i++ {
			encoded := encode(data[i*8:(i+1)*8], bs.WalshCode)
			signal = append(signal, encoded...)
		}
		fullSignals[b] = signal
		if timeSteps < len(signal) {
			timeSteps = len(signal)
		}
		fmt.Printf("Base Station %d: %d\n",bs.ID,bs.WalshCode)
    }

	// Send mixedsignals
	mixedSignals := make([]int, timeSteps)
	for x := 0; x < len(fullSignals); x++ {
		for y:=0; y < timeSteps; y++ {
			mixedSignals[y] += fullSignals[x][y]
		}
	}

	fmt.Println("Send mixed signals:",mixedSignals)

	// Discrete mixedsignals
	codeLength := len(mixedSignals) / walshCodeLength
	detectedWords := []string{}
	for x := 0; x < codeLength; x++ {
		mixedSignal := mixedSignals[x:x+8]
		detectedWords = mmseDetection(mixedSignal, baseStations)
	}

	fmt.Println("Summary of Detected Words:")
	for i, word := range detectedWords {
		fmt.Printf("Base Station %d: %s\n", baseStations[i].ID, word)
	}
}