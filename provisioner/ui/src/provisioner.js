import jQuery from "jquery"
import moment from "moment"

const hasProp = Object.prototype.hasOwnProperty

class Provisioner {
  constructor(urlDirectory) {
    this.urlDirectory = urlDirectory
    this.refreshRate  = false
    this.numDevices   = 0
    this.devices      = {}
    this.alerts       = []

    this.getDevices()
    this.startRefresh()
  }

  get refreshRate() { return this._refreshRate }
  set refreshRate(rate) {
    if (rate === false) {
      this._refreshRate = false
      this.stopRefresh()
      return
    }
    if (!typeof rate === "Number") {
      throw "refresh rate must be false or a number"
    }
    if (rate < 1000) {
      throw "refresh rate must be at least 1000ms"
    }
    this._refreshRate = rate
    this.startRefresh()
  }

  refreshHuman() {
    if (!this.refreshRate) {
      return "off"
    }
    let r = this.refreshRate/1000
    if (r < 60) {
      return `${r} sec`
    } else {
      return `${Math.round(r/60)} min`
    }
  }

  url(name, args={}) {
    let url = this.urlDirectory[name]
    if (!url) {
      throw `Named URL ${name} not found.`
    }

    (url.match(/\{[-\w]+\}/g) || []).forEach((match)=> {
      let el = match.substring(1, match.length-1) // "{foo}" → "foo"
      if (!hasProp.call(args, el)) {
        throw `Missing parameter ${el}`
      }
      url = url.replace(match, args[el])
    })

    return url
  }

  getDevices() {
    let promise = jQuery.getJSON(this.url("devices"))
    promise.done((data, _status, _xhr)=>{
      let devices = data
      devices.sort(function(a,b) {
        return +(a.mac_address < b.mac_address) || +(a.mac_address === b.mac_address) - 1
      })

      this.numDevices = devices.length
      this.devices = devices.reduce(function(memo, d) {
        memo[d.mac_address] = d
        return memo
      }, {})
    })
    promise.fail((xhr, status, error)=>{
      if (xhr.readyState === 0) {
        this.log("danger", "Could not connect to server.")
      } else {
        this.log("danger", `${status}: ${error}`)
      }
    })
    return
  }

  log(type, message) {
    this.alerts.unshift({ t: moment().unix(), style: `alert-${type}`, message: message })
    this.alerts.splice(5) // keep 5
  }

  startRefresh() {
    this.stopRefresh()
    if (this.refreshRate) {
      this._refreshID = setInterval(() => this.getDevices(), this.refreshRate)
    }
  }

  stopRefresh() {
    if (this._refreshID) {
      window.clearTimeout(this._refreshID)
      this._refreshID = null
    }
  }

  deviceAction(action, mac) {
    let dev = this.devices[mac]
    let url = null
    if (!dev) {
      this.log("danger", `Device ${mac} not found.`)
      return
    }
    if (["reboot", "provision", "upgrade"].indexOf(action) >= 0) {
      if (!(url = this.url(`${action}_device`, {mac: mac}))) {
        this.log("danger", `Don't know how to perform ${action} action: Missing route.`)
        return
      }
    } else {
      this.log("danger", `Unknown action: ${action}.`)
      return
    }
    if (action === "provision" && !dev.has_config) {
      this.log("danger", `No configuration found for device ${mac}.`)
      return
    }
    if (action === "upgrade" && !dev.can_upgrade) {
      this.log("danger", `No firmware upgrade found for device ${mac}.`)
      return
    }

    let promise = jQuery.ajax({ url: url, method: "POST" })
    promise.done((data, _status, _xhr) => {
      this.log(data.type, data.message)
    })
    promise.fail((xhr, status, error) => {
      let data = xhr.responseJSON
      if (data && data.type && data.message) {
        this.log(data.type, data.message)
      } else {
        this.log("danger", `Executing ${action} for device ${mac} failed (${status || error})`)
      }
    })
  }
}

export default Provisioner
