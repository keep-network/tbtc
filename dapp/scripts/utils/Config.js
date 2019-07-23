const fs = require('fs')

module.exports.readFile = function(path) {
  const configFile = fs.readFileSync(path, 'utf8')
  return JSON.parse(configFile)
}
