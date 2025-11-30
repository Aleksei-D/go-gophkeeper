package domain

import "context"

type UserStorage interface {
	GetCreds(ctx context.Context)
	CreateCreds(ctx context.Context)
	GetCard(ctx context.Context)
	GetCardsInfo(ctx context.Context)
	CreateCard(ctx context.Context)
	GetFile(ctx context.Context)
	SaveFile(ctx context.Context)
}
