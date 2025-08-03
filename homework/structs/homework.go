package main

import (
	"unsafe"
)

/*
Представим, что мы разрабатываем игру и наша задача упаковать данные игрока в структуру таким образом, чтобы ее размер был не более, чем 64 байта (*представим, что размер кэш линии на компьютерах будет 64 байта*).


**Данные пользователя:**

- __Имя пользователя__ \[0…42\] символов латиницы
  - **нельзя ссылаться на символы строки по указателю, например использовать тип данных string** (*нужно мапить символы строки в объект структуры, чтобы они находились рядом с другими данными*)
- __Координата по оси X__ \[-2_000_000_000…2_000_000_000\] значений
- __Координата по оси Y__ \[-2_000_000_000…2_000_000_000\] значений
- __Координата по оси Z__ \[-2_000_000_000…2_000_000_000\] значений
- __Золото__ \[0…2_000_000_000\] значений
- __Магическая сила (мана)__ \[0…1000\] значений
- __Здоровье__ \[0…1000\] значений
- __Уважение__ \[0…10\] значений
- __Сила__ \[0…10\] значений
- __Опыт__ \[0…10\] значений
- __Уровень__ \[0…10\] значений
- __Есть ли у игрока дом__ \[true/false\] значения
- __Есть ли у игрока оружие__ \[true/false\] значения
- __Есть ли у игрока семья__ \[true/false\] значения
- __Тип игрока__ \[строитель/кузнец/воин\] значения

*/

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type (
	Option func(*GamePerson)

	GamePerson struct {
		coordinates coordinates
		gold        uint32 // 0 - int32max
		name        [42]byte
		healthMana  [3]byte
		stats       [2]byte // respect | strength | exp | lvl
		flags       byte    // type---family|house|weapon
	}

	coordinates struct {
		x, y, z int32 // int32min - int32max
	}
)

var _ uintptr = unsafe.Sizeof(GamePerson{}) - 64
var _ uintptr = 64 - unsafe.Sizeof(GamePerson{})

func NewGamePerson(options ...Option) GamePerson {
	gp := &GamePerson{}

	for _, o := range options {
		o(gp)
	}

	return *gp
}

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		byteName := []byte(name)
		for i := range person.name {
			if i == len(byteName) {
				break
			}

			person.name[i] = byteName[i]
		}
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.coordinates = coordinates{x: int32(x), y: int32(y), z: int32(z)}
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		if gold < 0 {
			gold = 0
		}

		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		if mana < 0 {
			mana = 0
		}
		person.healthMana[1] |= byte((mana & 0x3F << 2))
		person.healthMana[2] = byte(mana >> 6)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		if health < 0 {
			health = 0
		}
		if health > 1000 {
			health = 1000
		}
		person.healthMana[0] = byte(health)
		person.healthMana[1] |= byte((health >> 8) & 0x03)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		v := byte(respect) & 0x0F
		person.stats[0] = person.stats[0] | v<<4
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		v := byte(strength) & 0x0F
		person.stats[0] = person.stats[0] | v
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		v := byte(experience) & 0x0F
		person.stats[1] = person.stats[1] | v<<4
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		v := byte(level) & 0x0F
		person.stats[1] = person.stats[1] | v
	}
}

func WithHouse() func(*GamePerson) {
	mask := byte(0b00000010)
	return func(person *GamePerson) {
		person.flags = person.flags | mask
	}
}

func WithGun() func(*GamePerson) {
	mask := byte(0b00000001)
	return func(person *GamePerson) {
		person.flags = person.flags | mask
	}
}

func WithFamily() func(*GamePerson) {
	mask := byte(0b00000100)
	return func(person *GamePerson) {
		person.flags = person.flags | mask
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.flags |= byte(personType) << 6
	}
}

func (p *GamePerson) Name() string {
	s := ""
	for _, b := range p.name {
		s += string(b)
	}

	return s
}

func (p *GamePerson) X() int {
	return int(p.coordinates.x)
}

func (p *GamePerson) Y() int {
	return int(p.coordinates.y)
}

func (p *GamePerson) Z() int {
	return int(p.coordinates.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {

	v := (int(p.healthMana[2])<<6 | int(p.healthMana[1]&0b11111100)>>2)
	return v
}

func (p *GamePerson) Health() int {
	v := (int(p.healthMana[1]&0b0000011)<<8 | int(p.healthMana[0]))
	return v
}

func (p *GamePerson) Respect() int {
	return int((p.stats[0] & 0b11110000) >> 4)
}

func (p *GamePerson) Strength() int {
	return int(p.stats[0] & 0x0F)
}

func (p *GamePerson) Experience() int {
	return int((p.stats[1] & 0b11110000) >> 4)
}

func (p *GamePerson) Level() int {
	return int(p.stats[1] & 0x0F)
}

func (p *GamePerson) HasHouse() bool {
	mask := byte(0b00000010)
	return p.flags&mask != 0
}

func (p *GamePerson) HasGun() bool {
	mask := byte(0b00000001)
	return p.flags&mask != 0
}

func (p *GamePerson) HasFamilty() bool {
	mask := byte(0b00000100)
	return p.flags&mask != 0
}

func (p *GamePerson) Type() int {
	return int(p.flags&(0b11<<6)) >> 6
}
