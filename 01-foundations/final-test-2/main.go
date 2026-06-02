package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// ============================================================================
// 🏆 FINAL FOUNDATIONS TEST 2: The Text Analyzer
// ============================================================================
// Your mission: Build a tool that analyzes text and finds the most used words!
//
// This challenge focuses specifically on what we learned in Module 5:
// ✅ Maps (Counting frequencies)
// ✅ Slices (Moving map data to a slice so we can sort it)
// ✅ Strings Package (Cleaning and splitting text)
// ✅ Structs (Grouping a word and its count together)
//
// INSTRUCTIONS: Follow the TODOs below to complete the program.
// RUN: go run main.go
// ============================================================================

// TODO 1: Create a struct called 'WordFreq' with two fields:
// - Word (string)
// - Count (int)

type WordFreq struct{
	Word string
	Count int
}
func main() {
	fmt.Println("=== 📚 Text Analyzer ===")
	fmt.Println("Type or paste some text. Press Enter on an empty line to finish.")

	scanner := bufio.NewScanner(os.Stdin)
	var allText string

	// TODO 2: Read multiple lines of input from the user.
	// - Use an infinite for loop with scanner.Scan()
	// - If the line is completely empty (input == ""), break the loop.
	// - Otherwise, add the line to 'allText' with a space at the end.
     for scanner.Scan() {
		input := scanner.Text()
		if input == "" {
			break
		}
		allText += input + " "
	 }
	// TODO 3: Clean the text!
	// - Convert 'allText' to lowercase using the 'strings' package.
	// - Replace commas ",", periods ".", and exclamation marks "!" with empty strings "".
	// (Hint: use strings.ToLower and strings.ReplaceAll)
    allText = strings.ToLower(allText)
	allText = strings.ReplaceAll(allText, ",", "")
	allText = strings.ReplaceAll(allText, ".", "")
	allText = strings.ReplaceAll(allText, "!", "")
	// TODO 4: Split the cleaned text into a slice of words.
	// (Hint: strings.Fields automatically splits by spaces and ignores extra spaces. It's much better than strings.Split for sentences!)
     words := strings.Fields(allText)
	 
	// TODO 5: Create a map to count how many times each word appears.
	// The key should be the word (string), and the value should be the count (int).
	// Loop through your slice of words from TODO 4 and update the map!
     wordCount := map[string]int{}
	 for _, word := range words{
		fmt.Printf("the word:%s", word)
		wordCount[word]++
	 }
	// TODO 6: We can't sort a map directly! So, let's move the data into a slice.
	// - Create an empty slice of your 'WordFreq' struct.
	// - Loop through your map using for...range.
	// - For each key/value pair, create a WordFreq struct and append it to the slice.
     var wordFreqs []WordFreq
	 for word, count := range wordCount{
		wordFreqs = append(wordFreqs, WordFreq{Word: word, Count: count})
	 }
	

	// TODO 7: Sort the slice by Count in descending order (highest count first).
	// Go provides a handy function for this. Uncomment and modify this block:
	// sort.Slice(yourSliceName, func(i, j int) bool {
	//     return yourSliceName[i].Count > yourSliceName[j].Count
	// })
    sort.Slice(wordFreqs, func(i, j int) bool {
		return wordFreqs[i].Count > wordFreqs[j].Count
	})
	// TODO 8: Print the top 3 most used words!
	// - Loop through the first 3 elements of your sorted slice.
	// - (Bonus points if you check the length of the slice first so you don't go out of bounds if the user typed less than 3 words!)
	// - Print them nicely, e.g., "1. 'the' (5 times)"
    for i, word := range wordFreqs{
		if i >= 3{
			break
		}
		fmt.Printf("%d. word:%s: %d\n", i+1, word.Word, word.Count)
	}
}
