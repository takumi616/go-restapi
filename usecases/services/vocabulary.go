package services

import (
	"context"
	"net/http"

	"github.com/takumi616/ielts-vocabularies-api/domains"
	"github.com/takumi616/ielts-vocabularies-api/usecases/ports"
)

type VocabService struct {
	Repo            domains.VocabRepository
	VocabOutputPort ports.VocabOutputPort
	ErrOutputPort   ports.ErrOutputPort
}

func (s *VocabService) AddNewVocabulary(ctx context.Context, vocab domains.Vocabulary, w http.ResponseWriter, err error) {
	//check if error is found in handlers
	if err != nil {
		s.ErrOutputPort.WriteErrResp(ctx, err, w, http.StatusBadRequest)
	}

	//Execute domain logic
	addedID, err := s.Repo.AddNewVocabulary(ctx, vocab)
	if err != nil {
		s.ErrOutputPort.WriteErrResp(ctx, err, w, http.StatusInternalServerError)
	}

	//Write http response
	s.VocabOutputPort.WriteVocabIdResp(ctx, addedID, w)
}