package validator

func HasArgAmount(args []string, amount int) bool {
	if len(args) < amount {
		return false
	}
	return true
}
