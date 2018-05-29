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
      return html`<div class="remote-file">
        - <a onclick=${download} href="#">${f.name}</a>
        ${downloadIndicator(state.downloads, p.name, f.name)}
      </div>`
    }
  }
}

function downloadIndicator(downloads, peer, file) {
  if (!downloads[`${peer}|${file}`]) return
  return html`<div class="progress-bar">
    <div class="progress-bar-inner" style="width: ${downloads[`${peer}|${file}`] * 100}%"></div>
  </div>`
}
