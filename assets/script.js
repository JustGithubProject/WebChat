const socket = new WebSocket("ws://localhost:8000/ws");

socket.addEventListener("open", () => {
    console.log("Connected to WebSocket server");
});

socket.addEventListener("message", (event) => {
    const msg = JSON.parse(event.data);
    const username = msg.username;
    const text = msg.text;

    const messagesList = document.getElementById("messages");
    const listItem = document.createElement("li");
    listItem.textContent = `${username}: ${text}`;
    messagesList.append(listItem);
});

socket.addEventListener("close", () => {
    console.log("Disconnected from WebSocket server");
});

document.getElementById("send").addEventListener("click", () => {
    const username = document.getElementById("username").value.trim();
    const message = document.getElementById("message").value.trim();

    if (username && message) {
        const msg = {
            username: username,
            text: message
        };

        socket.send(JSON.stringify(msg));
        document.getElementById("message").value = ''; 
    } else {
        alert("Enter your username and message");
    }
});
