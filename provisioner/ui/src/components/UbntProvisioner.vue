<template>
  <div class="ubnt-provisioner">
    <div class="pull-right small">
      <div class="btn-group btn-group-xs">
        <button type="button" class="btn btn-default dropdown-toggle" data-toggle="dropdown">
          Autorefresh: {{ provisioner.refreshHuman() }} <span class="caret"></span>
        </button>
        <button type="button" class="btn btn-default" v-on:click="provisioner.getDevices()">Refresh now</button>
        <ul class="dropdown-menu">
          <li><a href="#" v-on:click="provisioner.refreshRate = 1000">1 sec</a></li>
          <li><a href="#" v-on:click="provisioner.refreshRate = 15*1000">15 sec</a></li>
          <li><a href="#" v-on:click="provisioner.refreshRate = 60*1000">1 min</a></li>
          <li><a href="#" v-on:click="provisioner.refreshRate = 300*1000">5 min</a></li>
          <li><a href="#" v-on:click="provisioner.refreshRate = false">off</a></li>
        </ul>
      </div>
    </div>

    <h1 class="page-header">{{ title }}</h1>

    <div class="row">
      <section class="col-sm-6 col-md-8">
        <device-nav
            v-if="provisioner.numDevices"
            v-bind:currMac="currMac"
            v-bind:devices="provisioner.devices"
            v-on:device-navigate="onDeviceNavigate">
        </device-nav>

        <p class="alert alert-warning" v-else>
          No devices found. Please ensure that the host machine has an IP address in the <tt>169.254.0.0/16</tt> subnet.
        </p>

        <div class="error-log small" v-if="provisioner.alerts.length">
          <h3>Event log</h3>
          <div class="alert" v-for="a in provisioner.alerts" v-bind:class="a.style">
            <small class="pull-right">{{ a.t | timeAgo }}</small>
            {{ a.message }}
          </div>
        </div>
      </section>

      <aside class="col-sm-6 col-md-4">
        <device-view
            v-if="curr"
            v-bind:device="curr"
            v-on:device-action="onDeviceAction"
            v-on:close-view="onDeviceNavigate(null)">
        </device-view>
        <div class="alert alert-info" v-else>Please select a device.</div>
      </aside>
    </div>
  </div>
</template>

<script>
import moment     from "moment"
import filters    from "../utils/filters.js"
import DeviceView from "./DeviceView.vue"
import DeviceNav  from "./DeviceNav.vue"

export default {
  name: "UbntProvisioner",
  props: {
    provisioner: {
      type: Object,
      required: true
    },
  },
  data () {
    return {
      title: "Digineo Ubnt Provisioner",
      currMac: null
    }
  },
  computed: {
    curr: function() {
      if (this.currMac) {
        return this.provisioner.devices[this.currMac]
      }
    }
  },
  filters,
  methods: {
    onDeviceAction: function(name, mac) {
      this.provisioner.deviceAction(name, mac)
    },
    onDeviceNavigate: function(mac) {
      this.currMac = mac
    }
  },
  beforeMount() {
    this.provisioner.refreshRate = 15000
  },
  beforeDestroy () {
    this.provisioner.refreshRate = false
  },
  components: {
    DeviceView,
    DeviceNav
  }
}
</script>

<style scoped>
  .error-log .alert {
    margin-bottom: 1px;
    padding: 4px;
  }
</style>
