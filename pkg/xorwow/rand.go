package xorwow

// Code for xorwow taken from https://en.wikipedia.org/wiki/Xorshift

/* The state array must be initialized to not be all zero in the first four words */
type state struct {
	a, b, c, d uint32
	counter    uint32
}

func (s *state) xorwow() uint32 {
	/* Algorithm "xorwow" from p. 5 of Marsaglia, "Xorshift RNGs" */
	t := s.d

	v := s.a
	s.d = s.c
	s.c = s.b
	s.b = v

	t ^= t >> 2
	t ^= t << 1
	t ^= v ^ (v << 4)
	s.a = t

	s.counter += 362437
	return t + s.counter
}

// RNG : (Pseudo) Random Number Generator
type RNG interface {
	Rand() uint32
}

// NewRNG : create a new (Pseudo) Random Number Generator using the given seed
func NewRNG(seed uint32) RNG {
	init := &state{}
	init.a = seed
	for i := 0; i < 17; i++ {
		init.xorwow()
	}

	return &state{
		a:       init.xorwow(),
		b:       init.xorwow(),
		c:       init.xorwow(),
		d:       init.xorwow(),
		counter: init.xorwow(),
	}

}

func (s *state) Rand() uint32 {
	return s.xorwow()
}
