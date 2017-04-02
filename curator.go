package main

var MyCurator Curator

type Curator interface {
	// Init will be called on initialization, use this function to
	// initialize your curation module
	Init() error

	// OnPostAdded will be called when new posts are retrieved
	// from other peers, if this functions returns false, the
	// content will be rejected (e.g. in the case of spam) and not
	// stored by our node
	OnPostAdded(obj *Post, isWhitelabeled bool) bool

	// OnCommentAdded will be called when new comments are
	// retrieved from other peers, if this functions returns
	// false, the content will be rejected (e.g. in the case of
	// spam) and not stored by our node
	OnCommentAdded(obj *Comment, isWhitelabeled bool) bool

	// GetContent gives back an ordered array of post hashes of
	// suggested content by the curation module
	GetContent(params map[string]interface{}) []string

	// FlagContent marks or unmarks hashes as spam. True means
	// content is flagged as spam, false means content is not
	// flagged as spam
	FlagContent(hash string, isFlagged bool)

	// UpvoteContent is called when user upvotes a content
	UpvoteContent(hash string)

	// DownvoteContent is called when user downvotes a content
	DownvoteContent(hash string)

	// Close destruct curator module
	Close() error
}
