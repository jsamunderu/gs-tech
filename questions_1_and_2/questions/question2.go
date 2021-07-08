package questions

// IsFibNumber checks if a number is a member of the fibonacci series
func IsFibNumber(val int) bool {
	var n1, n2, n3 int = 1, 0, 0
	for {
		n3 = n1 + n2
		n1 = n2
		n2 = n3
		if n3 == val {
			return true
		}
		if val < n3 {
			return false
		}
	}
}
