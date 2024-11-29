package cronjob

import (
	"context"
	"fmt"
	advertise_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/robfig/cron/v3"
)

type crjPushToNewFeed struct {
	newFeedRepo advertise_repo.INewFeedRepository
}

func NewPushToNewFeedCronJob(
	newFeedRepo advertise_repo.INewFeedRepository,
) *crjPushToNewFeed {
	return &crjPushToNewFeed{
		newFeedRepo: newFeedRepo,
	}
}

func (crj *crjPushToNewFeed) Run() {
	ctx := context.Background()

	err := crj.newFeedRepo.CreateManyWithRandomUser(ctx, 100)
	if err != nil {
		fmt.Println("Error when pushing advertise to new feeds: ", err)
	}
}

func StartPushAdvertiseToNewFeedCronJob(
	newFeedRepo advertise_repo.INewFeedRepository,
) {
	c := cron.New()

	//_, err := c.AddFunc("@every 1m", func() {
	//	cronJob := NewPushToNewFeedCronJob(newFeedRepo)
	//	cronJob.Run()
	//})

	_, err := c.AddFunc("@daily", func() {
		cronJob := NewPushToNewFeedCronJob(newFeedRepo)
		cronJob.Run()
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	c.Start()
	fmt.Println("Push advertise cronjob started")

	select {}
}
