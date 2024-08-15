package cache

import (
	"dnevnik-rg.ru/internal/models"
	"sync"
)

type Cache struct {
	Pupils  Pupils
	Coaches Coaches
	Admins  Admins
}

type Pupils struct {
	Pupils map[int]*models.Pupil
	Mx     *sync.RWMutex
}

type Coaches struct {
	Coaches map[int]*models.Coach
	Mx      *sync.RWMutex
}

type Admins struct {
	Admins map[int]*models.Admin
	Mx     *sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{}
}

func (c *Cache) NewPupilsCache() IPupils {
	return &Pupils{
		Pupils: make(map[int]*models.Pupil),
		Mx:     &sync.RWMutex{},
	}
}

func (c *Cache) NewAdminsCache() IAdmin {
	return &Admins{
		Admins: make(map[int]*models.Admin),
		Mx:     &sync.RWMutex{},
	}
}

func (c *Cache) NewCoachesCache() ICoaches {
	return &Coaches{
		Coaches: make(map[int]*models.Coach),
		Mx:      &sync.RWMutex{},
	}
}
