package cronjob

import (
	"context"
	"fmt"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/cache"

	advertise_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/robfig/cron/v3"
)

type crjPushToNewFeed struct {
	newFeedRepo advertise_repo.INewFeedRepository
	postCache   cache.IPostCache
	cron        *cron.Cron
}

func NewPushToNewFeedCronJob(
	newFeedRepo advertise_repo.INewFeedRepository,
	postCache cache.IPostCache,
) *crjPushToNewFeed {
	crj := &crjPushToNewFeed{
		newFeedRepo: newFeedRepo,
		postCache:   postCache,
		cron:        cron.New(),
	}

	_, err := crj.cron.AddFunc("@every 6h", func() {
		go crj.Run()
	})

	if err != nil {
		fmt.Println(err)
		return crj
	}

	go crj.Run()

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

	if err = crj.postCache.DeleteAllPostCache(ctx); err != nil {
		fmt.Println("Error deleting all post cache: ", err)
	}
}

func (crj *crjPushToNewFeed) Stop() {
	crj.cron.Stop()
	fmt.Println("Cron job stopped")
}
