<h1>Websockets using Go standard library</h1>
This implementation is mainly for my own learnings to understand Go more and the underlying mechanics of websockets.
<br><br>
This is still WIP but if you'd like to give it a run, feel free to go through the steps below to get started

<h2>Prereqs</h2>
- Go
<br>
- Some client that can handle websocket requests. I recommend [websocat](https://github.com/vi/websocat) bc I like using the terminal.

<h2>Steps</h2>
- `cd` into the root of the folder
<br>
- start the Go server -> `go run .`
<br>
- in 2 separate terminals start websocat -> `websocat ws://localhost:3000/ws`
<br>
- start talking to yourself by typing messages in the websocat clients. 

<h2>Todo</h2>
[] create and fix project structure 
[] refactor code
[] create htmx client instead of websocat - cause i want to learn htmx
[] create rooms 
[] other things i cant think of right now
