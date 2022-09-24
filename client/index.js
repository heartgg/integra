console.log("Attempting Connection...");

const socket = new WebSocket(
  "ws://integri-scan.herokuapp.com/ws?roomID=1234&modality=XRAY"
);

socket.onopen = () => {
  console.log("Successfully Connected");
};

socket.onclose = (event) => {
  console.log("Socket Closed Connection: ", event);
};

socket.onerror = (error) => {
  console.log("Socket Error: ", error);
};

// handle messages received from server
socket.onmessage = (event) => {
  const msg = JSON.parse(event.data);
  console.log("Message received from server: ", msg);
  switch (msg.Type) {
    case 1:
      break;
    case 2:
      updatePatientInfo(msg);
      break;
    default:
      console.log(`Unknown message type ${msg.Type} received from server`);
      break;
  }
};

// update patient info displayed on page from message received from server
function updatePatientInfo(msg) {
  const infoList = document.getElementById("info-list");
  const examOpts = document.getElementById("exam-opts");
  const excludedOpts = document.getElementById("excluded-opts");

  let examCheckedCount = 0;

  infoList.innerHTML = "";
  examOpts.innerHTML = "";
  excludedOpts.innerHTML = "";
  examCheckedCount = 0;

  for (let key in msg.patient) {
    const li = document.createElement("li");
    li.innerHTML = `${key} : ${msg.patient[key]}`;
    infoList.appendChild(li);
  }
  let id = 0;
  for (let key in msg.exams) {
    const li = document.createElement("li");
    li.setAttribute("class", "list-group-item");
    const isSuggested = msg.exams[key] == 1 ? true : false;
    li.innerHTML = `
      <input
        class="form-check-input me-1"
        type="checkbox"
        value=""
        ${isSuggested ? "checked=true" : ""}
        id="checkbox-${id}"
      />
      <label class="form-check-label" for="checkbox-${id}"
        >${key}</label
      >`;
    if (isSuggested) {
      examOpts.appendChild(li);
      examCheckedCount = examCheckedCount + 1;
    } else {
      excludedOpts.appendChild(li);
    }
    id++;
  }
}
