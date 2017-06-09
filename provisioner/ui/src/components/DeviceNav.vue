<template>
  <div class="device-nav">
    <p class="input-group input-group-sm">
      <span class="input-group-addon">Filter (any field)</span>
      <input type="text" class="form-control" id="query" v-model="query">
      <span class="input-group-btn">
        <button class="btn btn-default" type="button"
            v-on:click="query = ''"
            v-bind:class="{ disabled: !query }">&times;</button>
      </span>
    </p>

    <table class="table table-condensed table-striped table-hover small">
      <thead>
        <tr>
          <th v-for="(header, key) in columns"
              v-on:click="sortBy(key)"
              v-bind:class="sortKey === key ? 'active' : ''">
            {{ header }}
            <span class="arrow" v-bind:class="sortOrders[key] > 0 ? 'asc' : 'desc'"></span>
          </th>
          <th></th>
        </tr>
      </thead>
      <tbody style="cursor: pointer">
        <tr v-for="dev in deviceList"
            v-on:click="navigateToDevice(dev)"
            v-bind:class="{ info: currMac === dev.mac_address }">
          <td v-for="key in keys"
              v-bind:class="sortKey === key ? 'active' : ''">{{ dev[key] }}</td>
          <td class="text-right">
            <span v-if="dev.has_config"  class="label label-info"   title="Configuration available">Cfg</span>
            <span v-if="dev.can_upgrade" class="label label-danger" title="Firmware upgrade available">Fw</span>
          </td>
        </tr>
      </tbody>
      <tfoot>
        <tr class="text-center">
          <td v-bind:colspan="keys.length + 1" class="small text-muted">
            <em>— click row for device details —</em>
          </td>
        </tr>
      </tfoot>
    </table>
  </div>
</template>

<script>
export default {
  name: "DeviceNav",
  data () {
    let columns = {
      platform:    "Platform",
      hostname:    "Hostname",
      ip_address:  "IP",
      mac_address: "MAC",
      status:      "Status",
    }
    let data = {
      columns:    columns,
      keys:       Object.keys(columns),
      query:      "",
      sortOrders: {},
    }

    data.sortKey = data.keys[0]
    data.keys.forEach(function(col) {
      return data.sortOrders[col] = 1
    })

    return data
  },
  computed: {
    deviceList: function() {
      let clean   = s => s.toLowerCase().replace(/[^\w\d]+/g, "")
      let sortKey = this.sortKey
      let query   = this.query && clean(this.query)
      let order   = this.sortOrders[sortKey] || 1
      let devices = Object.values(this.devices)

      if (query) {
        devices = devices.filter(function(dev) {
          return Object.keys(dev).some(function(key) {
            return clean(String(dev[key])).indexOf(query) >= 0
          })
        })
      }

      if (sortKey) {
        devices = devices.slice().sort(function(a,b) {
          a = a[sortKey]
          b = b[sortKey]
          return (a === b ? 0 : a > b ? 1 : -1) * order
        })
      }

      return devices
    }
  },
  props: {
    devices: { type: Object, required: true },
    currMac: { type: String }
  },
  methods: {
    navigateToDevice: function(dev) {
      this.$emit("device-navigate", dev.mac_address)
    },
    sortBy: function(key) {
      this.sortKey = key
      this.sortOrders[key] = this.sortOrders[key] * -1
    },
    arrowClass: function(key) {
      return this.sortOrders[key] > 0 ? "asc" : "desc"
    }
  }
}
</script>

<style scoped>
  .arrow {
    display:        inline-block;
    vertical-align: middle;
    width:          0;
    height:         0;
    margin-left:    5px;
    opacity:        0.66;
    border-left:    4px solid transparent;
    border-right:   4px solid transparent;
  }

  .arrow.asc  { border-bottom:  4px solid #333; }
  .arrow.desc { border-top:     4px solid #333; }
</style>
