package agent

func main() {
	calculators := NewFreeCalculators()
	calculators.RunCalculators() // Запуск вычислительных операций
	// TODO: Запуск демона RunDeamon(calculators)
}
