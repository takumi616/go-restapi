package database

import (
	"context"
	"log"

	"github.com/takumi616/ielts-vocabularies-api/adapters/dto"
	"gorm.io/gorm"
)

type Vocabulary struct {
	gorm.Model
	Title        string
	Definition   string
	Example      string
	PartOfSpeech string
	IsMemorized  string
}

// Insert new vocabulary data
func (g *Gorm) InsertNewVocabulary(ctx context.Context, vocabDto dto.VocabDto) (dto.VocabIdDto, error) {
	vocabulary := &Vocabulary{
		Title:        vocabDto.Title,
		Definition:   vocabDto.Definition,
		Example:      vocabDto.Example,
		PartOfSpeech: vocabDto.PartOfSpeech,
		IsMemorized:  vocabDto.IsMemorized,
	}

	result := g.Db.Create(vocabulary)
	vocabIdDto := dto.VocabIdDto{}
	if result.Error != nil {
		log.Printf("failed to insert new vocabulary: %v", result.Error)
		return vocabIdDto, result.Error
	}

	vocabIdDto.VocabularyID = vocabulary.ID

	return vocabIdDto, nil
}

func (g *Gorm) SelectVocabularyById(ctx context.Context, vocabularyID uint) (dto.VocabDto, error) {
	selected := Vocabulary{}
	result := g.Db.First(&selected, vocabularyID)
	vocabDto := dto.VocabDto{}
	if result.Error != nil {
		log.Printf("failed to select a vocabulary by ID: %v", result.Error)
		return vocabDto, result.Error
	}

	vocabDto.Title = selected.Title
	vocabDto.Definition = selected.Definition
	vocabDto.Example = selected.Example
	vocabDto.PartOfSpeech = selected.PartOfSpeech
	vocabDto.IsMemorized = selected.IsMemorized

	return vocabDto, nil
}
