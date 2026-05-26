package handlers

func isAllowedBanner(name string) bool {
	switch name {
	case "standard", "shadow", "thinkertoy":
		return true
	default:
		return false
	}
}
