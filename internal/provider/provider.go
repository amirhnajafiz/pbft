package provider

// Getter is interface for supporting GetAll method.
type Getter interface {
	GetAll() []*Item
}

// Adder is interface for supporting Add method.
type Adder interface {
	Add(item *Item)
}

// Item is a single news object.
type Item struct {
	Title string `json:"title"`
	Post  string `json:"post"`
}

// repo is our struct model.
// we don't need to have repo on public, we can
// create it as a private object and use its methods
// as public.
type repo struct {
	Items []*Item
}

// New creates a new Repo object
func New() *repo {
	return &repo{}
}

// Add function to support Adder interface.
func (r *repo) Add(item *Item) {
	r.Items = append(r.Items, item)
}

// GetAll function to support Getter interface.
func (r *repo) GetAll() []*Item {
	return r.Items
}
