package пакИнтерфейсы

//HitResult -- Used in mouse click events
type HitResult int

//ИСобытие -- интерфейс для событий
type ИСобытие interface{
	Type()
	Mod()
	Msg()
	X()
	Y()
	Err()
	Key()
	Ch()
	Width()
	Height()
	Target()
}