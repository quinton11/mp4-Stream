console.log("Hello World, Working");

const options = {
  iceServers: [
    {
      urls: ["stun:stun1.1.google.com:19302", "stun:stun2.1.google.com:19302"],
    },
  ],
  iceCandidatePoolSize: 10,
};

async function Init() {
  //Create peerconnection
  const peerconnection = new RTCPeerConnection(options);

  //create remote stream
  const remotestream = new MediaStream();
  document.getElementById("video-pion").srcObject = remotestream;

  peerconnection.onicecandidate = (event) => {
    if (event.candidate) {
      const JSON = {
        url: event.url,
        type: event.type,
        candidate: event.candidate,
        target: event.target,
      };
      console.log(JSON);
    }
  };

  //Add event listener for track
  peerconnection.ontrack = async (event) => {
    event.streams[0].getTracks().forEach((track) => {
      remotestream.addTrack(track);
    });
  };
  document.getElementById("button-sdp").addEventListener("click", async () => {
    const offer = await peerconnection.createOffer();

    peerconnection.setLocalDescription(offer);
    document.getElementById("sdp-1").value = JSON.stringify(offer);
    await postOffer("signal", offer);
    //When offer is created, send request to server and receive answer,
    //setremotedescription with answer and append it to textarea value
    //that should be all for signalling
  });
}

async function postOffer(endpoint, offer) {
  //Server url
  const url = "http://127.0.0.1:3000/signal";
  //const request = new Request()
  let response = await fetch(url, { method: "POST", body: offer });
  let data = await response.text()
  console.log(data);
}

Init();
