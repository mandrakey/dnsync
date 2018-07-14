package tools

func StringInSlice(s string, slice []string) bool {
    for _, item := range slice {
        if s == item {
            return true
        }
    }
    return false
}
