package taskmap

type TaskMap struct {
	item map[string]struct{}
}

func NewTaskMap() *TaskMap {
	return &TaskMap{item: make(map[string]struct{})}
}

func (m *TaskMap) Add(key string) bool {
	_, exists := m.item[key]
	if exists {
		return false // Key already exists
	}
	m.item[key] = struct{}{}
	return true // Key added successfully
}

func (m *TaskMap) Del(key string) {
	delete(m.item, key)
}

func (m *TaskMap) Get(key string) bool {
	_, exists := m.item[key]
	return exists
}

func (m *TaskMap) List() []string {
	keys := make([]string, 0, len(m.item))
	for key := range m.item {
		keys = append(keys, key)
	}
	return keys
}
