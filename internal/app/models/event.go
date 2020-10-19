package models

// Event is used to trigger event between engine.GameCommunication and engine.MainGameWindow
type Event interface {
	Dummy()
}

// AddCar is used to trigger a function to add car to joining player
type AddCar struct {}

// Dummy is here to satisfy Event interface requirement
func (addCar AddCar) Dummy() {}
