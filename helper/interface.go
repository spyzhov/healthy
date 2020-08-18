package helper

// Int is a function that return integer type of Value or Default
func Int(Value interface{}, Default int) int {
	switch Value := Value.(type) {
	case int:
		return Value
	case int8:
		return int(Value)
	case int16:
		return int(Value)
	case int32:
		return int(Value)
	case int64:
		return int(Value)
	case float32:
		return int(Value)
	case float64:
		return int(Value)
	case uint:
		return int(Value)
	case uint8:
		return int(Value)
	case uint16:
		return int(Value)
	case uint32:
		return int(Value)
	case uint64:
		return int(Value)
	default:
		return Default
	}
}
