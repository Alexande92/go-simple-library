package storage

// TODO make this storage not only book related, if needed

//type Entity Book | Book2

type Entity interface {
	getId() int64
	setId(s Table) int64
}

type Storage map[string]Table
type Table map[int]Entity

//var BookStorage = map[int]Book{}

func InitStorage() *Storage {
	return &Storage{"book": {}}
}

func FromTable(s Storage, name string) Table {
	return s[name]
}

func (s Table) Save(e Entity) {
	e.setId(s)
	(s)[len(s)+1] = e
}

func (s Table) GetAll() []Entity {
	entities := make([]Entity, 0, len(s))

	for _, v := range s {
		entities = append(entities, v)
	}
	return entities
}

func (s Table) GetById(id int64) (Entity, bool) {
	entity, ok := s[int(id)].(Entity)
	return entity, ok
}

func (s Table) Delete(id int64) {
	delete(s, int(id))
}

func (s Table) Update(e Entity) {
	s[int(e.getId())] = e
}
