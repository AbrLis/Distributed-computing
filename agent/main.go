package agent

func main() {
	calculators := NewFreeCalculators()
	RunCalculators()
	RunDeamon(calculators) // TODO: Запуск демона
}
