package cronjob

import (
	"context"
	"fmt"

	advertise_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/robfig/cron/v3"
)

type crjCheckExpiry struct {
	postRepo    advertise_repo.IPostRepository
	newFeedRepo advertise_repo.INewFeedRepository
}

func NewCheckExpiryCronJob(
	postRepo advertise_repo.IPostRepository,
	newFeedRepo advertise_repo.INewFeedRepository,
) *crjCheckExpiry {
	return &crjCheckExpiry{
		postRepo:    postRepo,
		newFeedRepo: newFeedRepo,
	}
}

func (crj *crjCheckExpiry) Run() {
	ctx := context.Background()

	err := crj.newFeedRepo.DeleteExpiredAdvertiseFromNewFeeds(ctx)
	if err != nil {
		fmt.Println("Error deleting expired advertisements")
	}

	err = crj.postRepo.UpdateExpiredAdvertisements(ctx)
	if err != nil {
		fmt.Println("Error updating expired advertisements")
	}
}

func StartCheckExpiryCronJob(
	postRepo advertise_repo.IPostRepository,
	newFeedRepo advertise_repo.INewFeedRepository,
) {
	c := cron.New()
	cronJob := NewCheckExpiryCronJob(postRepo, newFeedRepo)
	cronJob.Run()

	//_, err := c.AddFunc("@every 1m", func() {
	//	cronJob.Run()
	//})

	_, err := c.AddFunc("@daily", func() {
		cronJob.Run()
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	c.Start()
	fmt.Println("Check expiry of advertises cronjob started")

	select {}
}
