package main

func parseBool(v string) bool {
	switch v {
	case "1", "true", "TRUE":
		return true
	default:
		return false
	}
}
