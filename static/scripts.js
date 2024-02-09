const baseUrl = 'http://localhost:3000'
const sock = new WebSocket('ws://localhost:3000/ws')

sock.addEventListener("message", (event) => {
    console.log(event)
    messageData = JSON.parse(event.data)
    console.log(messageData)
    const chatbox = document.getElementById("chatbox")
    const newDiv = document.createElement('div')

    newDiv.textContent = messageData.message
    chatbox.appendChild(newDiv)
    chatbox.scrollTop = chatbox.scrollHeight
})

function sendMessage(event) {
    event.preventDefault()
    let input = document.getElementById("input")
    console.log(input.value)
    if (!input) return
    if (!sock.readyState) {
        console.error('socket not ready')

    }
    const messageData = {
        message: input.value
    }
    sock.send(JSON.stringify(messageData))
    input.value = ''
}
