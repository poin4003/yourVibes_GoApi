package cronjob

import (
	"context"
	"fmt"

	advertise_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/robfig/cron/v3"
)

type crjPushToNewFeed struct {
	newFeedRepo advertise_repo.INewFeedRepository
	cron        *cron.Cron
}

func NewPushToNewFeedCronJob(
	newFeedRepo advertise_repo.INewFeedRepository,
) *crjPushToNewFeed {
	crj := &crjPushToNewFeed{
		newFeedRepo: newFeedRepo,
		cron:        cron.New(),
	}

	_, err := crj.cron.AddFunc("@every 2m", func() {
		go crj.Run()
	})

	if err != nil {
		fmt.Println(err)
		return crj
	}

	go func() {
		crj.cron.Start()
		fmt.Println("Push advertise to new feed cronjob start")
	}()

	return crj
}

func (crj *crjPushToNewFeed) Run() {
	ctx := context.Background()

	err := crj.newFeedRepo.CreateManyWithRandomUser(ctx, 100)
	if err != nil {
		fmt.Println("Error when pushing advertise to new feeds: ", err)
	}
}

func (crj *crjPushToNewFeed) Stop() {
	crj.cron.Stop()
	fmt.Println("Cron job stopped")
}
