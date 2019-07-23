const createHash = require('create-hash')

exports.sha256 = function sha256(buffer) {
  return createHash('sha256')
    .update(buffer)
    .digest()
}

exports.hash256 = function hash256(buffer) {
  return this.sha256(this.sha256(buffer))
}
