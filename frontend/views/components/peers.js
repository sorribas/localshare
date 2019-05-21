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
  if (!state.peers.length) {
    return html`<center class="no-peers">No peers found</center>`
  }

  return html`
    <div>
      ${state.peers.map(renderPeer)}
    </div>
  `

  function renderPeer (p) {
    return html`<div>
      <div>${p.name}</div>
      ${fileTable()}
    </div>`

    function fileTable() {
      if (!p.files.length) return html`<div>- This peer hasn't shared any files yet.</div>`
      return html`<table class="remote-files"> ${p.files.map(renderRemoteFile)} </table>`
    }

    function renderRemoteFile (f) {
      var download = () => emit('download-file', {peer: p.name, file: f.name})
      return html`<tr>
        <td>- <a onclick=${download} href="#">${f.name}</a></td>
        <td>${downloadIndicator(state.downloads, p.name, f.name)}</td>
        <td>${savedFilePath(state.downloads, p.name, f.name)}</td>
      </tr>`
    }
  }
}

function downloadIndicator(downloads, peer, file) {
  if (!downloads[`${peer}|${file}`]) return
  var percent = Math.floor(downloads[`${peer}|${file}`].progress * 100)
  return html`<div class="progress-bar">
    <div class="progress-bar-inner" style="width: ${percent}%">${percent}%</div>
  </div>`
}

function savedFilePath (downloads, peer, file) {
  if (!downloads[`${peer}|${file}`]) return
  var filePath = downloads[`${peer}|${file}`].fileName
  var progress = downloads[`${peer}|${file}`].progress
  if (progress < 1) return `Downloading to ${filePath}`
  return `Downloaded to ${filePath}`
}
