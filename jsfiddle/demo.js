    /* eslint-env browser */

let pc = new RTCPeerConnection({
    iceServers: [
      {
        urls: 'stun:stun.l.google.com:19302'
      }
    ]
  })
  var log = msg => {
    document.getElementById('logs').innerHTML += msg + '<br>'
  }
  
  navigator.mediaDevices.getUserMedia({ video: true, audio: true })
    .then(stream => {
      stream.getTracks().forEach(track => pc.addTrack(track, stream));
      pc.createOffer().then(d => pc.setLocalDescription(d)).catch(log)
    }).catch(log)
    
  let sendChannel = pc.createDataChannel('foo')
  sendChannel.onclose = () => console.log('sendChannel has closed')
  sendChannel.onopen = () => console.log('sendChannel has opened')
  sendChannel.onmessage = e => log(`Message from DataChannel '${sendChannel.label}' payload '${e.data}'`)
  
  
  pc.oniceconnectionstatechange = e => log(pc.iceConnectionState)
  pc.onicecandidate = event => {
    if (event.candidate === null) {
      document.getElementById('localSessionDescription').value = btoa(JSON.stringify(pc.localDescription))
    }
  }
  pc.ontrack = function (event) {
    var el = document.createElement(event.track.kind)
    el.srcObject = event.streams[0]
    el.autoplay = true
    el.controls = true
  
    document.getElementById('remoteVideos').appendChild(el)
  }
  
  window.startSession = () => {
    let sd = document.getElementById('remoteSessionDescription').value
    if (sd === '') {
      return alert('Session Description must not be empty')
    }
  
    try {
      pc.setRemoteDescription(JSON.parse(atob(sd)))
    } catch (e) {
      alert(e)
    }
  }
  
  window.copySDP = () => {
    const browserSDP = document.getElementById('localSessionDescription')
  
    browserSDP.focus()
    browserSDP.select()
  
    try {
      const successful = document.execCommand('copy')
      const msg = successful ? 'successful' : 'unsuccessful'
      log('Copying SDP was ' + msg)
    } catch (err) {
      log('Unable to copy SDP ' + err)
    }
  }
  window.sendMessage = () => {
    let message = document.getElementById('message').value
    if (message === '') {
      return alert('Message must not be empty')
    }
  
    sendChannel.send(message)
  }