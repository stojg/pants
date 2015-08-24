package main

type Director struct {
	lastAction float64
}

func (director *Director) Update(w *World, frameTime, gameTime float64) {

	if gameTime-director.lastAction < 0.1 {
		return
	}

	//	if w.rand.Float32() > 0.9999999 {
	//		x := w.rand.Intn(w.worldMap.Width);
	//		y := w.rand.Intn(w.worldMap.Height)
	//		w.worldMap.Base().Place(x, y, 5)
	//		w.networkManager.BroadcastMap(w.worldMap.Compress())
	//	}

	if len(entityManager.sprites) >= 25 {
		return
	}
	director.lastAction = gameTime
	w.entityManager.NewEntity(
		w.rand.Float64()*1000,
		w.rand.Float64()*1000,
		"assets/basics/arrow.png",
	)
}
