package initialize

import (
	"github.com/poin4003/yourVibes_GoApi/global"
	advertise_cronjob "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/cronjob"
	post_repo_impl "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/post/repo_impl"
)

func InitCronJob() {
	db := global.Pdb

	postRepo := post_repo_impl.NewPostRepositoryImplement(db)
	newFeedRepo := post_repo_impl.NewNewFeedRepositoryImplement(db)

	go advertise_cronjob.StartPushAdvertiseToNewFeedCronJob(newFeedRepo)
	go advertise_cronjob.StartCheckExpiryCronJob(postRepo, newFeedRepo)
}
