package provider

// Getter is interface for supporting GetAll method.
type Getter interface {
	GetAll() []*Item
}

// Adder is interface for supporting Add method.
type Adder interface {
	Add(item *Item)
}

type Item struct {
	Title string `json:"title"`
	Post  string `json:"post"`
}

// Repo is our struct model.
type Repo struct {
	Items []*Item
}

func New() *Repo {
	return &Repo{
		Items: []*Item{},
	}
}

// Add function to support Adder interface.
func (r *Repo) Add(item *Item) {
	r.Items = append(r.Items, item)
}

// GetAll function to support Getter interface.
func (r *Repo) GetAll() []*Item {
	return r.Items
}
