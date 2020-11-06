package internal

import "../framework"

func Ping(ctx framework.Context) {
	ctx.Reply("pong")
}