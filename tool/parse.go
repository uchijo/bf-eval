package tool

func Parse(input []uint8) []uint8 {
	retval := []uint8{}
	for _, c := range input {
		switch c {
		case '>', '<', '+', '-', '.', ',', '[', ']':
			retval = append(retval, c)
		}
	}
	return retval
}
