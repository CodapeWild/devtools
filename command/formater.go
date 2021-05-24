package command

type OutputFormater func([]byte) ([]byte, error)
