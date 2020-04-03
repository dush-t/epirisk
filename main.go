package main

import (
	"log"
	"net/http"

	"github.com/dush-t/epirisk/api"
	"github.com/dush-t/epirisk/middleware"
)

func main() {
	config := InitializeApp()

	// TODO Find a way to make this middleware chaining less ugly
	http.Handle("/sign_in", middleware.LogReq(api.SignInHandler(config.DBConn, config.Bus)))
	http.Handle("/sign_up", middleware.LogReq(api.SignUpHandler(config.DBConn, config.Bus)))
	http.Handle("/met_user", middleware.LogReq(middleware.Auth(config.DBConn, api.MetUserHandler(config.DBConn, config.Bus))))
	http.Handle("/update_self_healthstatus", middleware.LogReq(middleware.Auth(config.DBConn, api.UpdateSelfHealthStatus(config.DBConn, config.Bus))))
	http.Handle("/get_contact_summary", middleware.LogReq(middleware.Auth(config.DBConn, api.GetContactSummary(config.DBConn, config.Bus))))

	log.Println("HTTP server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
