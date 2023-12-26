package cache

import "dnevnik-rg.ru/internal/models"

type IAdmin interface {
	ReadAll() map[int]*models.Admin
	ReadById(key int) (*models.Admin, bool)
	WriteAdmin(admin models.Admin)
	RemoveAdmin(key int)
	WritingSession(list []models.Admin)
}

func (a *Admins) ReadAll() map[int]*models.Admin {
	a.Mx.RLock()
	defer a.Mx.RUnlock()
	return a.Admins
}

func (a *Admins) ReadById(key int) (*models.Admin, bool) {
	a.Mx.RLock()
	defer a.Mx.RUnlock()
	if _, ok := a.Admins[key]; !ok {
		return &models.Admin{}, false
	}
	return a.Admins[key], true
}

func (a *Admins) WriteAdmin(admin models.Admin) {
	a.Mx.Lock()
	defer a.Mx.Unlock()
	a.Admins[admin.Key] = &admin
}

func (a *Admins) RemoveAdmin(key int) {
	a.Mx.Lock()
	defer a.Mx.Unlock()
	delete(a.Admins, key)
}

func (a *Admins) WritingSession(list []models.Admin) {
	a.Mx.Lock()
	defer a.Mx.Unlock()
	for _, admin := range list {
		a.Admins[admin.Key] = &admin
	}
}
