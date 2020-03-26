package main

import (
	"log"
	"net/http"

	"github.com/dush-t/epirisk/api"
	"github.com/dush-t/epirisk/init"
	"github.com/dush-t/epirisk/middleware"
)

func main() {
	config := init.InitializeApp()

	// TODO Find a way to make this middleware chaining less ugly
	http.Handle("/sign_in", middleware.LogReq(api.SignInHandler(config)))
	http.Handle("/sign_up", middleware.LogReq(api.SignUpHandler(config)))
	http.Handle("/met_user", middleware.LogReq(middleware.Auth(config, api.MetUserHandler(config))))
	http.Handle("/update_self_healthstatus", middleware.LogReq(middleware.Auth(config, api.UpdateSelfHealthStatus(config))))
	http.Handle("/get_contact_summary", middleware.LogReq(middleware.Auth(config, api.GetContactSummary(config))))

	log.Println("HTTP server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
