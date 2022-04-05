package backoff

type Attempts interface {
	SameAttempts | DifferentAttempts

	Next(int) int
}

type DifferentAttempts []int

func (attempts DifferentAttempts) Next(i int) int {
	if len(attempts) == 0 {
		return 0
	}

	if len(attempts) <= i {
		return attempts[len(attempts)-1]
	}

	return attempts[i]
}

type SameAttempts int

func (attempts SameAttempts) Next(int) int {
	return int(attempts)
}
