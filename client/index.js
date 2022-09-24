// import { connect } from "./api";
const socket = new WebSocket("ws://integri-scan.herokuapp.com/ws?roomID=1234&modality=XRAY");

// connect to websocket and handle messages
const infoList = document.getElementById("info-list")
const examOpts = document.getElementById("exam-opts")
const excludedOpts = document.getElementById("excluded-opts");

let examCheckedCount = 0;

socket.addEventListener('message', (event) => {
  console.log('Message from server ', event.data);

  infoList.innerHTML = '';
  examOpts.innerHTML = '';
  excludedOpts.innerHTML = '';
  // FIXME: Change back later to not hard-coded data
  // const msg = JSON.parse(event.data);
  const msg = {
    "patient": {
      "ID": 1293811,
      "Name": "Brad J", 
      "Birthdate": "07/18/2002", 
      "Sex": "Male", 
      "Diagnosis": "Skin Cancer"
    },
    "exams": {
      "Angiography": 0,
      "Arthrography": 0,
      "Bone Density Scan": 1,
      "Bone XRAY": 0,
      "Chest XRAY": 1,
      "Crystogram": 0,
      "Fluoroscopy": 1,
      "Mammography": 0,
      "Myelography": 0,
      "Skull Radiography": 0,
      "Virtual Colonoscopy": 1
    }
  };

  for (let key in msg.patient) {
    const li = document.createElement("li");
    li.innerHTML = `${key} : ${msg.patient[key]}`;
    infoList.appendChild(li);
  }
  for (let key in msg.exams) {
    const li = document.createElement("li");
    li.setAttribute("class", "list-group-item");
    const isSuggested = (msg.exams[key] == 1 ? true : false);
    li.innerHTML = `
      <input
        class="form-check-input me-1"
        type="checkbox"
        value=""
        ${(isSuggested ? 'checked=true' : '')}
        id="firstCheckbox"
      />
      <label class="form-check-label" for="firstCheckbox"
        >${key}</label
      >`;
      if (isSuggested) {
        examOpts.appendChild(li);
        examCheckedCount = examCheckedCount + 1;
      }
      else {
        excludedOpts.appendChild(li);
      }
    
    }
  switch (msg.Type) {
    case 1:
      console.log(msg)
      break;
    case 2:
      console.log("Message to be processed", msg.Body)
      break;
    default:
      // skip
      break;
  }
});

