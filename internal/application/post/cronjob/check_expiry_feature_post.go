package cronjob

import (
	"context"
	"fmt"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/cache"
	advertise_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/robfig/cron/v3"
)

type crjCheckExpiryFeaturePost struct {
	newFeedRepo advertise_repo.INewFeedRepository
	postCache   cache.IPostCache
	cron        *cron.Cron
}

func NewCheckExpiryFeaturePostCronJob(
	newFeedRepo advertise_repo.INewFeedRepository,
	postCache cache.IPostCache,
) *crjCheckExpiryFeaturePost {
	crj := &crjCheckExpiryFeaturePost{
		newFeedRepo: newFeedRepo,
		postCache:   postCache,
		cron:        cron.New(),
	}

	_, err := crj.cron.AddFunc("@every 5h", func() {
		go crj.Run()
	})

	if err != nil {
		fmt.Println("Error when check expiry feature post: ", err)
		return crj
	}

	go crj.Run()

	go func() {
		crj.cron.Start()
		fmt.Println("Check feature post cronjob started")
	}()

	return crj
}

func (crj *crjCheckExpiryFeaturePost) Run() {
	ctx := context.Background()

	err := crj.newFeedRepo.DeleteExpiredFeaturedPostsFromNewFeeds(ctx)
	if err != nil {
		fmt.Println(err)
	}

	if err = crj.postCache.DeleteAllPostCache(ctx); err != nil {
		fmt.Println("failed to delete all post cache: ", err)
	}
}

func (crj *crjCheckExpiryFeaturePost) Stop() {
	crj.cron.Stop()
	fmt.Println("Cronjob check feature expiry stopped")
}
