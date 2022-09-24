import { connect } from "./api";

// connect to websocket and handle messages
const infoList = document.findById("info-list")
const examOpts = document.findById("exam-opts")
connect((msg) => {
  for (const key of msg.patientData) {
    const li = document.createElement("li");
    li.innerHTML = `${key} : ${msg[key]}`;
    infoList.appendChild(li);
  }
  for (const key of msg.examSuggestions) {
    const li = document.createElement("li");
    li.setAttribute("class", "list-group-item");
    li.innerHTML = `
      <input
        class="form-check-input me-1"
        type="checkbox"
        value=""
        id="firstCheckbox"
      />
      <label class="form-check-label" for="firstCheckbox"
        >First checkbox</label
      >`;
      examOpts.appendChild(li);
      }
});
