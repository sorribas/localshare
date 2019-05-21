var html = require('choo/html')

module.exports = shareComponent 

function shareComponent (state, emit) {
  return html`
    <div class="section">
      <h2>Shared files</h2>
      ${fileList(state, emit)}
      <button onclick=${clickBtn} class="btn">Add file</button>
    </div>
  `

  function clickBtn () {
    emit('add-file')
  }
}

function fileList (state, emit) {
  if (!state.sharedFiles.length) {
    return html`<center class="no-peers">No files shared</center>`
  }

  return html`
    <div>
      ${state.sharedFiles.map(renderFile)}
    </div>
  `

  function renderFile (file) {
    return html`<div>${file.name} <a href="#" onclick=${copy}>(Copy link)</a></div>`

    function copy() {
      copytoclipboard(`${state.serverUrl}/${file.name}`)
    }
  }
}
