package frodo

// Little function to convert type "func(w http.ResponseWriter, r *http.Request)" to Frodo.HandleFunc
func makeHandleFunc(h HandleFunc) HandleFunc {
	// Log.Debug("converting func(http.ResponseWriter, *Request) to Frodo.HandleFunc")
	return h
}

// func(http.ResponseWriter, *http.Request, Params)
func makeRouterHandleFunc(rh HTTPRouterHandleFunc) HTTPRouterHandleFunc {
	// Log.Debug("converting func(http.ResponseWriter, *http.Request, Params) to httprouter.Handle")
	return rh
}
