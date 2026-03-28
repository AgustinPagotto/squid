package validator

func HasArgAmount(args []string, amount int) bool {
	return len(args) < amount
}
