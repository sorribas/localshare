var html = require('choo/html')

module.exports = peersComponent

function peersComponent (state, emit) {
  return html`
    <div>
      <h2>Peers</h2>
      ${peerList(state, emit)}
    </div>
  `
}

function peerList(state, emit) {
  return html`
    <div>
      ${state.peers.map(renderPeer)}
    </div>
  `

  function renderPeer (p) {
    return html`<div>
      <div>${p.name}</div>
      <div>${p.files.map(renderRemoteFile)}</div>
    </div>`

    function renderRemoteFile (f) {
      var download = () => emit('download-file', {peer: p.name, file: f.name})
      return html`<div class="remote-file">- <a onclick=${download} href="#">${f.name}</a></div>`
    }
  }
}
