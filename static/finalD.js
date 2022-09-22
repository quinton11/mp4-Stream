//Creating web socket
let ws = new WebSocket("ws://localhost:3000"); //connection url

let peerconnection;
let remotestream = new MediaStream();
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
  console.log(event.candidate);
  if (event.candidate === null) {
    console.log("ICE gathered!");

    //updating offer textarea
    document.getElementById("sdp-1").value = JSON.stringify(
      peerconnection.localDescription
    );
  }
};

//Event handler for a remote track
peerconnection.ontrack = (event) => {
  console.log("Ontrack event");
  //RTC Media Stream Event
  console.log(event);
  //Media Stream Object
  console.log(event.streams[0]);
  //MediaStream Track
  console.log(event.track);
  //remotestream.addTrack(event.track);
  remotestream = event.streams[0]

  document.getElementById("videos-main").srcObject = remotestream

  //element to hold video
  /* var elem = document.createElement(event.track.kind);
  elem.srcObject = event.streams[0];
  elem.autoplay = true;
  elem.controls = true;

  var par = document.getElementById("video-pion");

  par.appendChild(elem); */

  return false;
};

peerconnection.oniceconnectionstatechange = (event) => {
  console.log("Ice Event State:");
  console.log(peerconnection.iceConnectionState);
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
    //peerconnection.setRemoteDescription(answer);
  } else {
    //ICE candidates
    console.log("Peer ICE candidates: " + JSON.stringify(answer));
    //peerconnection.addIceCandidate(answer);
  }
};

//creating Starting offer
//stored in textarea
peerconnection.createOffer({ iceRestart: true }).then((offer) => {
  peerconnection.setLocalDescription(offer);

  //store offer in textarea
  document.getElementById("sdp-1").value = JSON.stringify(offer);
});

let triggerStream = async () => {
  const url = "http://localhost:3000/streamupdate";

  let response = await fetch(url, {
    method: "POST",
    headers: { Accept: "application/json", "Content-Type": "application/json" },
    body: JSON.stringify({ stream: "start" }),
  });
  let data = await response.json();
};

let submitOffer = async () => {
  const url = "http://localhost:3000/signal";
  const offr = document.getElementById("sdp-1").value;

  let response = await fetch(url, {
    method: "POST",
    headers: { Accept: "application/json", "Content-Type": "application/json" },
    body: offr,
  });
  let data = await response.json();
  peerconnection.setRemoteDescription(data);
};

//Add event listener for button
//when clicked, start stream
document.getElementById("button-sdp").addEventListener("click", triggerStream);
document.getElementById("button-sdpO").addEventListener("click", submitOffer);
