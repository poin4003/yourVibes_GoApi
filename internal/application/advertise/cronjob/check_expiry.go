package cronjob

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"

	advertise_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
)

type crjCheckExpiry struct {
	postRepo    advertise_repo.IPostRepository
	newFeedRepo advertise_repo.INewFeedRepository
	cron        *cron.Cron
}

func NewCheckExpiryCronJob(
	postRepo advertise_repo.IPostRepository,
	newFeedRepo advertise_repo.INewFeedRepository,
) *crjCheckExpiry {
	crj := &crjCheckExpiry{
		postRepo:    postRepo,
		newFeedRepo: newFeedRepo,
		cron:        cron.New(),
	}

	_, err := crj.cron.AddFunc("@every 12h", func() {
		go crj.Run()
	})
	if err != nil {
		fmt.Println(err)
		return crj
	}

	go crj.Run()

	go func() {
		crj.cron.Start()
		fmt.Println("Check expiry of advertises cronjob started")
	}()

	return crj
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

func (crj *crjCheckExpiry) Stop() {
	crj.cron.Stop()
	fmt.Println("Check expiry Cronjob stopped")
}
