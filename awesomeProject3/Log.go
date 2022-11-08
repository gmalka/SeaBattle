package main

const stackSize = 9

type log struct {
	str   [stackSize]string
	index byte
}

func (s *log) add(newElement string) {
	if s.index < stackSize {
		s.str[s.index] = newElement
		s.index++
	} else {
		for i := 1; i < stackSize; i++ {
			s.str[i] = s.str[i-1]
		}
		s.str[stackSize-1] = newElement
	}
}

func (s *log) get(index int) string {
	if index < stackSize {
		return s.str[index]
	} else {
		return ""
	}
}
