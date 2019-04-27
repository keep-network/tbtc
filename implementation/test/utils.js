

async function deploySystem(deploy_list) {
  deployed = {}  // name: contract object
  linkable = {}  // name: linkable address
  for (let i in deploy_list) {
    await deploy_list[i].contract.link(linkable)
    contract = await deploy_list[i].contract.new()
    linkable[deploy_list[i].name] = contract.address
    deployed[deploy_list[i].name] = contract
  }
  return deployed
}

module.exports = {
  deploySystem: deploySystem
}
