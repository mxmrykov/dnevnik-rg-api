package cache

import "dnevnik-rg.ru/internal/models"

type IPupils interface {
	ReadAll() map[int]*models.Pupil
	ReadById(key int) (*models.Pupil, bool)
	WritePupil(pupil models.Pupil)
	RemovePupil(key int)
	WritingSession(list []models.Pupil)
}

func (p *Pupils) ReadAll() map[int]*models.Pupil {
	p.Mx.RLock()
	defer p.Mx.RUnlock()
	return p.Pupils
}

func (p *Pupils) ReadById(key int) (*models.Pupil, bool) {
	p.Mx.RLock()
	defer p.Mx.RUnlock()
	if _, ok := p.Pupils[key]; !ok {
		return &models.Pupil{}, false
	}
	return p.Pupils[key], true
}

func (p *Pupils) WritePupil(pupil models.Pupil) {
	p.Mx.Lock()
	defer p.Mx.Unlock()
	p.Pupils[pupil.Key] = &pupil
}

func (p *Pupils) RemovePupil(key int) {
	p.Mx.Lock()
	defer p.Mx.Unlock()
	delete(p.Pupils, key)
}

func (p *Pupils) WritingSession(list []models.Pupil) {
	p.Mx.Lock()
	defer p.Mx.Unlock()
	for _, pupil := range list {
		pupilCopy := pupil
		p.Pupils[pupil.Key] = &pupilCopy
	}
}
