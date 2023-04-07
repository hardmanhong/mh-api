package goods

type Item struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	BasePrice float32 `json:"basePrice"`
}

type Goods interface {
	Create(goods *Item) error
	Update(goods *Item) error
	Delete(id int) error
	GetItem(id int) (Item, error)
	GetList() ([]Item, error)
}

func (g *Item) Create() error {
	return nil
}
func (goods *Item) Update() error {
	return nil
}
func (goods *Item) Delete() error {
	return nil
}
func GetItem(id int) (*Item, error) {
	return nil, nil
}
func GetList() (*[]Item, error) {
	return nil, nil
}
