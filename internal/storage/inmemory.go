package storage

import "errors"

// TODO make this storage not only book related, if needed

type Entity interface {
	getId() int
	setId(s Table) int
}

type Storage map[string]Table
type Table map[int]Entity

func InitStorage() Storage {
	return Storage{"book": {}}
}

func FromTable(s Storage, name string) Table {
	return s[name]
}

func (s Table) Save(e Entity) {
	id := e.setId(s)
	(s)[id] = e
}

func (s Table) GetAll() []Entity {
	entities := make([]Entity, 0, len(s))

	for _, v := range s {
		entities = append(entities, v)
	}
	return entities
}

func (s Table) GetById(id int) (Entity, error) {
	entity, ok := s[id].(Entity)

	if !ok {
		return nil, errors.New("book not found")
	}
	return entity, nil
}

func (s Table) Delete(id int) {
	delete(s, id)
}

func (s Table) Update(e Entity) {
	s[e.getId()] = e
}
