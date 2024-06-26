package service

import (
	"back/internal/repository"
	"back/internal/schemas"
)

type Authorization interface {
	SignUp(userSchema *schemas.CreateUserReq) (*schemas.CreateUserResp, error)
	SignIn(userSchema *schemas.SignInReq) (*schemas.SignInResp, error)
	GetProfile(userID int) (*schemas.GetProfileResp, error)
	UpdateUsername(usernameSchema *schemas.UpdateUsernameReq) (*schemas.UpdateUsernameResp, error)
	UpdatePassword(passwordSchema *schemas.UpdatePasswordReq) error
}

type Collection interface {
	CreateCollection(collectionSchema *schemas.CreateCollectionReq, userID int) (*schemas.CreateCollectionResp, error)
	UpdateCollection(collectionSchema *schemas.UpdateCollectionReq) (*schemas.UpdateCollectionResp, error)
	RemoveCollection(collectionSchema *schemas.RemoveCollectionReq) error
	GetAllCollections(userID int) (*schemas.AllCollectionsResp, error)
	TrainCards(req *schemas.TrainSchemaReq) (*schemas.TrainSchemaResp, error)
}

type Card interface {
	CreateCard(cardSchema *schemas.CreateCardReq, collectionID int) (*schemas.CreateCardResp, error)
	UpdateCard(cardSchema *schemas.UpdateCardReq) (*schemas.UpdateCardResp, error)
	RemoveCard(cardSchema *schemas.RemoveCardReq) error
	GetCardsByCollectionID(collectionID int) (*schemas.GetCardByCollectionIDResp, error)
}

type Service struct {
	Authorization
	Collection
	Card
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.UserRepository),
		Collection:    NewCollectionService(repos.CollectionRepository),
		Card:          NewCardService(repos.CardRepository),
	}
}
