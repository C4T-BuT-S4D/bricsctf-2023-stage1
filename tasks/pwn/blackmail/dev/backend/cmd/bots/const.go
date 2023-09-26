package main

import (
	"cbs.dev/brics/droidchat/internal/app"
	"math/rand"
)

var thePassword = "06b6468e247ead496a2f"

var usernames = []string{
	"Mercedes60",
	"Peter_Stracke37",
	"Marcella.Rohan92",
	"Eloy_Torp",
	"Rosamond.Trantow18",
	"Kelli29",
	"Domenico_Stiedemann",
	"Berneice_Bins",
	"Arlene_Breitenberg88",
	"Rowland_Stanton",
	"Desiree_Pollich80",
	"Joey.Abshire",
	"Ivory.Blick",
	"Lavada71",
}

func text(t string) app.Message {
	return app.Message{
		Type: app.MessageTypeText,
		Text: &t,
	}
}

var responses = []app.Message{
	text("LOL!"),
	text("That's so cool!"),
	text("What is the meaning of Life?"),
	text("Gotcha."),
	{Type: app.MessageTypeSticker, Sticker: &app.Sticker{Id: "meh"}},
	text("Hello!"),
	{Type: app.MessageTypeSticker, Sticker: &app.Sticker{Id: "yay"}},
	text("Greetings!"),
	text("I like your attitude."),
	{Type: app.MessageTypeSticker, Sticker: &app.Sticker{Id: "yay"}},
	text("Thank you for sharing that."),
	text("Wait, no way?!"),
	text("Classic..."),
	{Type: app.MessageTypeSticker, Sticker: &app.Sticker{Id: "owo"}},
	text("Uhh... what is that?"),
}

func randomMsg() app.Message {
	return responses[rand.Intn(len(responses))]
}
