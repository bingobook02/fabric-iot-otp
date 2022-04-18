package main

// func main1() {
// 	// iot env
// 	os.Setenv("DISCOVERY_AS_LOCALHOST", "true")

// 	r := mux.NewRouter()

// 	r.HandleFunc("/add-device", addDevice)

// 	// Solves Cross Origin Access Issue
// 	c := cors.New(cors.Options{
// 		AllowedOrigins: []string{"http://localhost:4200"},
// 	})
// 	handler := c.Handler(r)

// 	srv := &http.Server{
// 		Handler: handler,
// 		Addr:    ":" + os.Getenv("PORT"),
// 	}

// 	log.Fatal(srv.ListenAndServe())
// }

// func addDevice(w http.ResponseWriter, r *http.Request) {
// 	var data = struct {
// 		Title string `json:"title"`
// 	}{
// 		Title: "Golang + Angular Starter Kit",
// 	}

// 	jsonBytes, err := utils.StructToJSON(data)
// 	if err != nil {
// 		fmt.Print(err)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(jsonBytes)
// 	return
// }
