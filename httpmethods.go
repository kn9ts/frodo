package frodo

// Methods type is used in Match method to get all methods user wants to apply
// that will help in invoking the related handler
type Methods []string

// MethodsAllowed -- HTTP Methods/Verbs allowed
var MethodsAllowed = Methods{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
