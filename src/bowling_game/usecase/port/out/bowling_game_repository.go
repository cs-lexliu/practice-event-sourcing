package out

import (
	"github.com/cs-lexliu/practice-event-sourcing/src/bowling_game/entity"
	core "github.com/cs-lexliu/practice-event-sourcing/src/core/usecase"
)

type BowlingGameRepository core.Repository[*entity.BowlingGame]
