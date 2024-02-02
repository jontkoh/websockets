# Websockets using Go standard library
This implementation is mainly for my own learnings to understand Go more and the underlying mechanics of websockets.

This is still WIP but if you'd like to give it a run, feel free to go through the steps below to get started

## Prereqs
- Go
- Some client that can handle websocket requests. I recommend [websocat](https://github.com/vi/websocat) bc I like using the terminal. 

## Steps
- `cd` into the root of the folder
- start the Go server -> `go run .`
- in 2 separate terminals start websocat -> `websocat ws://localhost:3000/ws`
- start talking to yourself by typing messages in the websocat clients. 

## Todo
- [ ] create and fix project structure 
- [ ] refactor code
- [ ] create htmx client instead of websocat - cause i want to learn htmx
- [ ] create rooms 
- [ ] other things i cant think of right now
