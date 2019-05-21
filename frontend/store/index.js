module.exports = store

function store (state, emitter) {
  state.peers = []
  state.sharedFiles = [{name: 'test.txt'}]
  state.downloads = {}

  emitter.on('download-file', function (data) {
    localshare.download(data.peer, data.file)
  })

  emitter.on('add-file', function () {
    localshare.chooseFile()
  })

  window.update = function () {
    state.peers = localshare.data.peers
    state.sharedFiles = localshare.data.files
    state.downloads = localshare.data.downloads
    state.serverUrl = localshare.data.serverUrl
    emitter.emit('render')
  }
}
