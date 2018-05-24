var html = require('choo/html')
var shareComponent = require('./components/share')
var peersComponent = require('./components/peers')

module.exports = mainView

function mainView (state, emit) {
  return html`
    <div>
      <div class="header"><h1>LocalShare</h1></div>
      <div class="container">
        ${shareComponent(state, emit)}
        <hr />
        ${peersComponent(state, emit)}
      </div>
    </div>
  `
}
