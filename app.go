package main

import (
	"context"
	"fmt"
	"kltRPA/models"
	"kltRPA/utils"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// omo登录
func (a *App) Login(username, password string) string {
	omo := utils.NewOmoIntegrate("omo.kelote.com", "klt_omo", username, password)

	success, err := omo.Login()
	if success {
		return "Login successful!"
	} else {
		if err != nil {
			return fmt.Sprintf("Login failed: %v", err)
		} else {
			return "Invalid username or password."
		}
	}
}

func (a *App) RunRPA() {
	models.RunRPA()
}
