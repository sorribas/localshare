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
  return html`
    <div>
      ${state.sharedFiles.map(renderFile)}
    </div>
  `

  function renderFile (file) {
    return html`<div>${file.name}</div>`
  }
}
