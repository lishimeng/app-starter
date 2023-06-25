package masker

func Phone(input string) (s string) {
	if len(input) != 11 {
		return input
	}

	s = input[:3] + "****" + input[len(input)-4:]
	return
}

func Email(input string) (s string) {

	return
}
