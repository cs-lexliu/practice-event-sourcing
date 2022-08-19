package out

import (
	"github.com/cs-lexliu/practice-event-sourcing/internal/ddd/cusecase"
	"github.com/cs-lexliu/practice-event-sourcing/pkg/domain/bowling_game/entity"
)

type BowlingGameRepository cusecase.Repository[*entity.BowlingGame]
