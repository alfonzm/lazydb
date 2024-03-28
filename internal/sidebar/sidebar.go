package sidebar

type Sidebar struct {}

func New() *Sidebar {
  return &Sidebar{}
}

func (s *Sidebar) Draw() string {
  return "rendered sidebar"
}
