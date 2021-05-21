module.exports = async function () {
  const networkID = await web3.eth.net.getId()
  console.log(networkID)
  process.exit(0)
}
