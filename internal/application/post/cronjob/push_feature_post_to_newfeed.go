package cronjob

import (
	"context"
	"fmt"
	advertise_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/robfig/cron/v3"
)

type crjPushFeaturePostToNewFeed struct {
	newFeedRepo advertise_repo.INewFeedRepository
	cron        *cron.Cron
}

func NewPushFeaturePostToNewFeedCronJob(
	newFeedRepo advertise_repo.INewFeedRepository,
) *crjPushFeaturePostToNewFeed {
	crj := &crjPushFeaturePostToNewFeed{
		newFeedRepo: newFeedRepo,
		cron:        cron.New(),
	}

	_, err := crj.cron.AddFunc("@every 2h", func() {
		go crj.Run()
	})

	if err != nil {
		fmt.Println("Error when pushing feature post to new feeds: ", err)
		return crj
	}

	go func() {
		crj.cron.Start()
		fmt.Println("Push feature post cron job started")
	}()

	return crj
}

func (crj *crjPushFeaturePostToNewFeed) Run() {
	ctx := context.Background()

	err := crj.newFeedRepo.CreateManyFeaturedPosts(ctx, 100)
	if err != nil {
		fmt.Println("Error when pushing advertise to new feeds: ", err)
	}
}

func (crj *crjPushFeaturePostToNewFeed) Stop() {
	crj.cron.Stop()
	fmt.Println("Cron job stopped")
}
