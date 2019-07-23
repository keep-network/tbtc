module.exports.fromHex = function(hex) {
  return Buffer.from(hex, 'hex')
}

module.exports.toHex = function(bytes) {
  const buffer = Buffer.from(bytes)
  return buffer.toString('hex')
}
