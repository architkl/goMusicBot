package internal

import "../framework"

func Username(ctx framework.Context) {
	ctx.Reply("Your username is " + ctx.Message.Author.Username)
}