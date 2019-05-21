var choo = require('choo')
var devtools = require('choo-devtools')
var store = require('./store')

var app = choo()
app.use(devtools())
app.use(store)
app.route('*', require('./views/main'))
app.mount('#app')
