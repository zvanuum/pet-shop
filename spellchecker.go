package main

import (
	"log"
	"strings"

	"github.com/sajari/fuzzy"
)

type ISpellChecker interface {
	NewSpellChecker()
	CountMisspelledWords(string) int
	IsWordMisspelled(string) bool
}

type SpellChecker struct {
	Model *fuzzy.Model
}

func (checker *SpellChecker) NewSpellChecker() {
	checker.Model = fuzzy.NewModel()
	checker.Model.SetThreshold(1)
	checker.Model.SetDepth(2)
	checker.Model.Train(fuzzy.SampleEnglish())
	log.Printf("Finished building spell checker\n")
}

func (checker *SpellChecker) CountMisspelledWords(tweet string) int {
	var misspelledCount int
	words := strings.Split(tweet, " ")

	for _, word := range words {
		if  !strings.Contains(word, "http") && checker.IsWordMisspelled(word) {
			misspelledCount++
		}
	}

	return misspelledCount
}

func (checker *SpellChecker) IsWordMisspelled(word string) bool {
	// for some reason the spell checker doesn't recognize I or a
	isAorI := word == "a" || word == "i" || word == "I"

	return isAorI || word != checker.Model.SpellCheck(word)
}
