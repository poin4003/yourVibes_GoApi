package cronjob

import (
	"context"
	"fmt"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/cache"
	advertise_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/robfig/cron/v3"
)

type crjPushFeaturePostToNewFeed struct {
	newFeedRepo advertise_repo.INewFeedRepository
	postCache   cache.IPostCache
	cron        *cron.Cron
}

func NewPushFeaturePostToNewFeedCronJob(
	newFeedRepo advertise_repo.INewFeedRepository,
	postCache cache.IPostCache,
) *crjPushFeaturePostToNewFeed {
	crj := &crjPushFeaturePostToNewFeed{
		newFeedRepo: newFeedRepo,
		postCache:   postCache,
		cron:        cron.New(),
	}

	_, err := crj.cron.AddFunc("@every 5h", func() {
		go crj.Run()
	})

	if err != nil {
		fmt.Println("Error when pushing feature post to new feeds: ", err)
		return crj
	}

	go crj.Run()

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

	if err = crj.postCache.DeleteAllPostCache(ctx); err != nil {
		fmt.Println("Error when deleting all post cache: ", err)
	}
}

func (crj *crjPushFeaturePostToNewFeed) Stop() {
	crj.cron.Stop()
	fmt.Println("Cron job stopped")
}
