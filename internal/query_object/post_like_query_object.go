package query_object

type PostLikeQueryObject struct {
	PostID string `form:"post_id,omitempty"`
	Limit  int    `form:"limit,omitempty"`
	Page   int    `form:"page,omitempty"`
}
