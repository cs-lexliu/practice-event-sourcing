package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cs-lexliu/practice-event-sourcing/internal/ddd/adpater"
	"github.com/cs-lexliu/practice-event-sourcing/pkg/domain/bowling_game/entity"
	"github.com/cs-lexliu/practice-event-sourcing/pkg/domain/bowling_game/usecase"
	"github.com/cs-lexliu/practice-event-sourcing/pkg/domain/bowling_game/usecase/port/in"
	"github.com/google/uuid"
)

func main() {

	reader := bufio.NewReader(os.Stdin)

	eventBus := adpater.NewDomainEventBus()
	repository := adpater.NewEventRepositoryAdapter[*entity.BowlingGame](eventBus)

	var currentID string

loop:
	for {
		fmt.Println("What's your command: create, roll, exit")
		cmd := getInput(reader)
		switch cmd {
		case "create":
			currentID = uuid.New().String()
			u := usecase.NewCreateBowlingGameUseCase(repository)
			input := in.CreateBowlingGameInput{BowlingGameID: currentID}
			if err := u.Execute(context.Background(), input); err != nil {
				fmt.Println(err)
				break loop
			}
		case "roll":
			fmt.Println("How many pins you hit? 0 - 10")
			hitStr := getInput(reader)
			hit, err := strconv.Atoi(hitStr)
			if err != nil {
				fmt.Println(err)
				break loop
			}
			u := usecase.NewRollOneBallUseCase(repository)
			input := in.RollOneBallInput{BowlingGameID: currentID, Hit: hit}
			if err := u.Execute(context.Background(), input); err != nil {
				fmt.Println(err)
				break loop
			}
		case "exit":
			break loop
		}
		b, err := repository.FindByID(context.Background(), currentID)
		if err != nil {
			fmt.Println(err)
			break loop
		}
		fmt.Printf("Your scores is %d!!!\n", b.Scores())
	}
}

func getInput(reader *bufio.Reader) string {
	cmd, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return strings.TrimSuffix(cmd, "\n")
}
