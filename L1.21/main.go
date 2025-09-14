package main

import (
	"fmt"
)

// старый чайник
type Kettle struct {
	isActive bool
}

func (k *Kettle) kettelStop() {
	k.isActive = false
}

func (k *Kettle) kettleStart() {
	k.isActive = true
}

// старая печь
type Stove struct {
	isBurning bool
}

func (s *Stove) extinguish() {
	s.isBurning = false
}

func (s *Stove) ignite() {
	s.isBurning = true
}

// интерфейс для нового типа устройств
type SmartDevice interface {
	switchOn()
	switchOff()
}

// лампа, которая новая и реализует новый интерфейс
type Lamp struct {
	isActive bool
}

func (l *Lamp) switchOn() {
	l.isActive = true
}

func (l *Lamp) switchOff() {
	l.isActive = false
}

// ниже пошли адаптеры под старые объекты
type KettleAdapter struct {
	kettle *Kettle
}

func (ka *KettleAdapter) switchOn() {
	ka.kettle.kettleStart()
}

func (ka *KettleAdapter) switchOff() {
	ka.kettle.kettelStop()
}

type StoveAdapter struct {
	stove *Stove
}

func (sa *StoveAdapter) switchOn() {
	sa.stove.ignite()
}

func (sa *StoveAdapter) switchOff() {
	sa.stove.extinguish()
}

func main() {

	//создали объекты
	kettle := &Kettle{}
	stove := &Stove{}
	lamp := &Lamp{}

	kettleAdapter := &KettleAdapter{kettle: kettle}
	stoveAdapter := &StoveAdapter{stove: stove}
	//ну а лампе адаптер не нужен, она и так реализовывает новый интерфейс

	devices := []SmartDevice{kettleAdapter, stoveAdapter, lamp}
	for _, device := range devices {
		device.switchOn()
	}

	//проверим состояния
	fmt.Println("Lamp active:", lamp.isActive)
	fmt.Println("Kettle active:", kettle.isActive)
	fmt.Println("Stove burning:", stove.isBurning)
}
