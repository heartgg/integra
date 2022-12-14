const DEV = false;

console.log("Attempting Connection...");

function connectSocket(url) {
  const socket = new WebSocket(url);
  
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
        updatePatientInfo(JSON.parse(msg.Body));
        break;
      default:
        console.log(`Unknown message type ${msg.Type} received from server`);
        break;
    }
  };
  return socket;
}

 var currentSocket = connectSocket((DEV ? "ws://localhost:3000" : "wss://integri-scan.herokuapp.com") + "/ws?roomID=1234&modality=XRAY");

// update patient info displayed on page from message received from server
function updatePatientInfo(data) {
  console.log(data);
  const infoList = document.getElementById("info-list");
  const examOpts = document.getElementById("exam-opts");
  const excludedOpts = document.getElementById("excluded-opts");
  const loadedDataDiv = document.getElementById("loaded-data-container")
  const noDataDiv = document.getElementById("no-data-label")
  const collapseUnsuggestedBtn = document.getElementById("collapseButton");
  
  let examCheckedCount = 0;

  infoList.innerHTML = "";
  examOpts.innerHTML = "";
  excludedOpts.innerHTML = "";
  
  loadedDataDiv.classList.remove("not-visible");
  noDataDiv.classList.add("not-visible"); 
  
  for (let key in data.Patient) {
    const tr = document.createElement("tr");
    const tdLeft = document.createElement("td");
    const tdRight = document.createElement("td");
    tdLeft.innerHTML = key[0].toUpperCase() + key.substring(1);
    if (key == "sex") {
      tdRight.innerHTML = data.Patient[key] ? "Male" : "Female";
    } else if (key == "birthdate") {
      tdRight.innerHTML = data.Patient[key].split("T")[0];
    } else {
      tdRight.innerHTML = data.Patient[key];
    }
    // li.innerHTML = `${key} : ${data.Patient[key]}`;
    tr.appendChild(tdLeft);
    tr.appendChild(tdRight);
    infoList.appendChild(tr);
  }

  let id = 0;
  for (let key in data.Exams) {
    const li = document.createElement("li");
    li.setAttribute("class", "list-group-item");
    const isSuggested = data.Exams[key] == 1 ? true : false;
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
      checkDisableButton("confirm-btn", examCheckedCount);
    } else {
      excludedOpts.appendChild(li);
    }
    let input = li.querySelector("input");
    if (input != null) {
      input.addEventListener("input", (event) => {
        if (input.checked == true) {
          // Then the user just checked the box
          examCheckedCount++;
        } else {
          // Then the user just unchecked the box
          examCheckedCount--;
        }
        checkDisableButton("confirm-btn", examCheckedCount);
      });
    }
    id++;
  }


}

document.getElementById("collapseButton").click(); 


function checkDisableButton(buttonId, num) {
  const btn = document.getElementById(buttonId);
  if (num <= 0) {
    btn.disabled = true;
  } else {
    btn.disabled = false;
  }
}

const btn = document.getElementById("collapseButton");
var collapsed = false;
btn.addEventListener("click", function handleClick() {
  const initialText = "Other Exam Options";
  collapsed = !collapsed;
  console.log(btn.innerHTML);
  if (collapsed) {
    btn.innerHTML = `<i class="bi bi-arrow-down"></i>&nbspShow`;
  } else {
    btn.innerHTML = `<i class="bi bi-arrow-up"></i>&nbspHide`;
  }
});

const modality = ["IE", "Fluoro", "XRAY", "CT", "IR", "MRI", "US", "Dexa", "NucMed"]

for (const item of modality) {
  const li = document.createElement("li");
  li.setAttribute("class", "dropdown-item");
  li.innerHTML = `
  <a class="dropdown-item" href="#">${item}</a>`;
  let input = li.querySelector("a");
  input.addEventListener("click", () =>{
    currentSocket.close();
    currentSocket = connectSocket((DEV ? "ws://localhost:3000" : "wss://integri-scan.herokuapp.com") + `/ws?roomID=1234&modality=${input.text}`);
    headerText = document.getElementById("header-text").innerText = input.innerText + " Workstation";
    showInfo(false);
  });
  document.getElementById("modality-dropdown").appendChild(li);
}

document.getElementById("confirm-btn").addEventListener("click", () => {
    showInfo(false);
    runSuccess();
});

function runSuccess() {
  $(".success-checkmark").show();
  $(".check-icon").show();
  setTimeout(function () {
    $(".check-icon").hide();
    $(".success-checkmark").hide();
  }, 1500);
}

$(".check-icon").hide();
$(".success-checkmark").hide();

function showInfo (boolean) {
  const loadedDataDiv = document.getElementById("loaded-data-container");
  const noDataDiv = document.getElementById("no-data-label");
  if (boolean) {
    loadedDataDiv.classList.remove("not-visible");
    noDataDiv.classList.add("not-visible"); 
  } else {
    loadedDataDiv.classList.add("not-visible");
    noDataDiv.classList.remove("not-visible"); 
  }
  collapsed = false;
}