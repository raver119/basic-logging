package logging

type Field struct {
	Key   string
	Value interface{}
}

// F is a helper function for creating a Field
func F(key string, value interface{}) Field {
	return Field{key, value}
}
