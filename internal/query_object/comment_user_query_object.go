package query_object

type CommentQueryObject struct {
	PostId   string `form:"post_id" binding:"required"`
	ParentId string `form:"parent_id,omitempty"`
	Limit    int    `form:"limit,omitempty"`
	Page     int    `form:"page,omitempty"`
}
