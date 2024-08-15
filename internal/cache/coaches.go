package cache

import "dnevnik-rg.ru/internal/models"

type ICoaches interface {
	ReadAll() map[int]*models.Coach
	ReadById(key int) (*models.Coach, bool)
	WriteCoach(coach models.Coach)
	RemoveCoach(key int)
	WritingSession(list []models.Coach)
}

func (a *Coaches) ReadAll() map[int]*models.Coach {
	a.Mx.RLock()
	defer a.Mx.RUnlock()
	return a.Coaches
}

func (a *Coaches) ReadById(key int) (*models.Coach, bool) {
	a.Mx.RLock()
	defer a.Mx.RUnlock()
	if _, ok := a.Coaches[key]; !ok {
		return &models.Coach{}, false
	}
	return a.Coaches[key], true
}

func (a *Coaches) WriteCoach(coach models.Coach) {
	a.Mx.Lock()
	defer a.Mx.Unlock()
	a.Coaches[coach.Key] = &coach
}

func (a *Coaches) RemoveCoach(key int) {
	a.Mx.Lock()
	defer a.Mx.Unlock()
	delete(a.Coaches, key)
}

func (a *Coaches) WritingSession(list []models.Coach) {
	a.Mx.Lock()
	defer a.Mx.Unlock()
	for _, coach := range list {
		coachCopy := coach
		a.Coaches[coach.Key] = &coachCopy
	}
}
