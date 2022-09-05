//Creating web socket
let ws = new WebSocket("ws://127.0.0.1:3000"); //connection url

let peerconnection;
//A check if websocket is open
ws.onopen = () => {
  alert("Socket Connected!");
};

//STUN servers for configuring our peer connection
const options = {
  iceServers: [
    {
      urls: ["stun:stun1.1.google.com:19302", "stun:stun2.1.google.com:19302"],
    },
  ],
  iceCandidatePoolSize: 10,
};

//peerconnection object
peerconnection = new RTCPeerConnection(options);

//Event for when we receive ICE candidates from STUN servers
peerconnection.onicecandidate = (event) => {
  //Candidates trickle in
  //When all candidates have been gathered
  //a final null candidate is passed
  if (event.candidate === null) {
    console.log("ICE gathered!");

    //updating offer textarea
    document.getElementById("sdp-1").value = JSON.stringify(
      peerconnection.localDescription
    );

    //when all ice candidates have been gathered then send offer
    ws.send(JSON.stringify(peerconnection.localDescription));
  }
};

//Event handler for a remote track
peerconnection.ontrack = (event) => {
  console.log("Ontrack event");

  //element to hold video
  var elem = document.createElement(event.track.kind);
  elem.srcObject = event.streams[0];
  elem.autoplay = true;
  elem.controls = true;

  document.getElementById("video-pion").appendChild(elem);
  return false;
};

//creating transceiver for peerconnection
//For application, we'll be receiving
//video track from server
peerconnection.addTransceiver("video", { direction: "recvonly" });

//when we receive response from websocket
ws.onmessage = (event) => {
  let answer = JSON.parse(event.data);

  //Check if answer was received
  //ICE candidates are being sent
  //in addition
  if (answer["type"] === "answer") {
    console.log(answer);
    peerconnection.setRemoteDescription(answer);
  } else {
    //ICE candidates
    console.log("Peer ICE candidates: " + JSON.stringify(answer));
  }
};

//A setTimeout to display connection state after a while
let delay = 10000;
setTimeout(() => {
  console.log(peerconnection.iceConnectionState);
}, delay);

//creating Starting offer
//stored in textarea
peerconnection.createOffer({ iceRestart: true }).then((offer) => {
  //wait for websocket to go live
  setTimeout(async () => {
    await peerconnection.setLocalDescription(offer);

    //store offer in textarea
    document.getElementById("sdp-1").value = JSON.stringify(offer);
  }, 3000);
});

let triggerStream = async () => {
  const url = "http://127.0.0.1:3000/stream";

  let response = await fetch(url, {
    method: "POST",
    headers: { Accept: "application/json", "Content-Type": "application/json" },
    body: JSON.stringify({ stream: "start" }),
  });
  let data = await response.json();
};

//Add event listener for button
//when clicked, start stream
document.getElementById("button-sdp").addEventListener("click", triggerStream);
