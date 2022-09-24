const socket = new WebSocket("ws://integri-scan.herokuapp.com/ws?roomID=1234&modality=XRAY");

const connect = callback => {
  console.log("Attempting Connection...");

  socket.onopen = () => {
    console.log("Successfully Connected");
  };

  socket.onmessage = msg => {
    console.log(msg);
    callback(msg);
  };

  socket.onclose = event => {
    console.log("Socket Closed Connection: ", event);
  };

  socket.onerror = error => {
    console.log("Socket Error: ", error);
  }
};

const sendMsg = (msg) => {
  console.log("Sending Message: ", msg);
  socket.send(msg);
};

export { connect, sendMsg };