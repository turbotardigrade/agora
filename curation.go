package main

var MyCurator = DefaultCurator{}

type Curator interface {
	Init() error
	OnContentAdded(obj *IPFSObj) bool
	GetContent(params map[string]interface{})
	FlagContent(isFlagged bool)
	UpvoteContent(hash string)
	DownvoteContent(hash string)
}

type DefaultCurator struct{}

func (c *DefaultCurator) Init() error {
	return nil
}

func (c *DefaultCurator) OnContentAdded(obj *IPFSObj) bool {
	return true
}

func (c *DefaultCurator) GetContent(params map[string]interface{}) []string {
	// @TODO return string of hashes
	return []string{}
}

func (c *DefaultCurator) FlagContent(isFlagged bool) {

}

func (c *DefaultCurator) UpvoteContent(hash string) {

}

func (c *DefaultCurator) DownvoteContent(hash string) {

}
