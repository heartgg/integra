import { connect } from "./api";

// connect to websocket and handle messages
const infoList = document.findById("info-list")
connect((msg) => {
  for (const key of msg) {
    const li = document.createElement("li");
    li.innerHTML = `${key} : ${msg[key]}`;
    infoList.appendChild(li);
  }
});