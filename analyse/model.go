package analyse

type UserAnalyseItem struct {
	Email string
	Name  map[string]bool
	Plus  int
	Minus int
}

func New(email string) *UserAnalyseItem {
	return &UserAnalyseItem{
		Email: email,
		Name:  make(map[string]bool),
		Minus: 0,
		Plus:  0,
	}
}

func (uai *UserAnalyseItem) AddRecord(name string, plus, minus int) {
	uai.Plus += plus
	uai.Minus += minus
	uai.Name[name] = true
}
