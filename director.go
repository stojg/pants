package main

type Director struct {
	lastAction float64
}

func (director *Director) Update(w *World, frameTime, gameTime float64) {

	if gameTime-director.lastAction < 0.1 {
		return
	}

	if len(entityManager.entities) >= 400 {
		return
	}

	entType := ENT_ENT1
	rand := w.rand.Float64()
	if rand > 0.9 {
		entType = ENT_FOOD
	}

	director.lastAction = gameTime
	w.entityManager.NewEntity(
		w.rand.Float64()*1000,
		w.rand.Float64()*1000,
		entType,
	)
}
