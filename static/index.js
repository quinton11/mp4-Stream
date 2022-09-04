console.log("Init");
console.log("Hello World, Working");
//PeerConnection and remotestream variables
let peerconnection;
let remotestream;
let ws = new WebSocket("ws://127.0.0.1:3000");

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

//Creating a  PeerConnection
peerconnection = new RTCPeerConnection(options);

//our remotestream
remotestream = new MediaStream();

//peerconnection event handlers
peerconnection.onicecandidate = (event) => {
  if (event.candidate) {
    //console.log(event.candidate);
    //peerconnection.addIceCandidate(event.candidate);
    ws.send(peerconnection.LocalDescription);
  }
};

peerconnection.iceConnectionState = (event) => {
  console.log("Ice connection state" + peerconnection.iceConnectionState);
};

//listening for tracks
peerconnection.ontrack = (event) => {
  /* event.streams[0].getTracks().forEach((track) => {
    remotestream.addTrack(track);
  }); */
  remotestream.addTrack(event.track);
  console.log("Ontrack event!");
  console.log(event.track)
  console.log(event.receiver)
  console.log(event.transceiver)

  //console.log(event.streams[0]);
  document.getElementById("video-pion").srcObject = remotestream;
  return false;
};

//document.getElementById("video-pion").srcObject = remotestream;
//*creating offer
//first we create transceivers, for this case we use one video media
peerconnection.addTransceiver("video", { direction: "sendrecv" });

ws.onmessage = async (event) => {
  await peerconnection.setRemoteDescription(event.data);
};

//offer
peerconnection.createOffer({ iceRestart: "true" }).then(async (offer) => {
  await peerconnection.setLocalDescription(offer);

  //store offer in textarea element
  document.getElementById("sdp-1").value = JSON.stringify(offer);
});

//an event to get offer from textarea and send it as a request to server
//response is received and used to set remote description
async function sendOffer() {
  let obj = document.getElementById("sdp-1");
  //read value from sdp-1 element
  let offr = JSON.parse(obj.value);

  //send offer to server
  let answer = await postOffer(offr);

  //setting remote description
  await peerconnection.setRemoteDescription(answer);
  //append answer to document element for viewing

  //obj.value += obj.value + "\\n Answer: \\n" + JSON.stringify(answer);
}
document.getElementById("button-sdp").addEventListener("click", sendOffer);

async function postOffer(offer) {
  //Server url

  const url = "http://127.0.0.1:3000/signal";
  //const request = new Request()
  let response = await fetch(url, {
    method: "POST",
    headers: { Accept: "application/json", "Content-Type": "application/json" },
    body: JSON.stringify(offer),
  });
  let data = await response.json();
  //console.log(data);
  return data;
}
