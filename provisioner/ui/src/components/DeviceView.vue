<template>
  <div class="device-view">
    <h2>
    </h2>

    <div class="panel panel-info">
      <div class="panel-heading">
        <h4 class="panel-title">
          Device Details
          <span class="badge">{{ device.status }}</span>
        </h4>
      </div>
      <div class="panel-body">
        <dl class="dl-horizontal">
          <dt>Hostname</dt>
          <dd>{{ device.hostname }}</dd>
          <dt>MAC address</dt>
          <dd><code>{{ device.mac_address }}</code></dd>
          <dt>Model</dt>
          <dd>{{ device.model || "n/a" }}</dd>
          <dt>Platform</dt>
          <dd>{{ device.platform }}</dd>
          <dt>primary IP address</dt>
          <dd><code>{{ device.ip_address }}</code></dd>
          <dt>other IP addresses</dt>
          <dd v-for="(ips, mac) in device.ip_addresses">
            <code>{{mac}}</code>
            <ul>
              <li v-for="ip in ips"><code>{{ip}}</code></li>
            </ul>
          </dd>
        </dl>
      </div>
      <div class="panel-body">
        <dl class="dl-horizontal">
          <dt>Firmware</dt>
          <dd><tt>{{device.firmware}}</tt> <span class="label label-danger" v-if="device.can_upgrade">!</span></dd>

          <dt>Configuration</dt>
          <dd>{{ device.has_config ? "available" : "not available" }}</dd>

          <dt>first seen</dt>
          <dd>{{ device.first_seen_at | fmtDate }}</dd>
          <dt>last seen</dt>
          <dd>{{ device.last_seen_at | fmtDate }}</dd>

          <dt>up since</dt>
          <dd>{{ device.up_since | fmtDate }}</dd>
        </dl>
      </div>
      <div class="panel-body">
        <dl class="dl-horizontal">
          <dt>ESSID</dt>
          <dd>{{ device.essid }}</dd>
          <dt>Wireless mode</dt>
          <dd>{{ device.wireless_mode }}</dd>
        </dl>
      </div>
      <div class="panel-footer" v-if="device.status !== 'idle'">
        <p>The device is busy ({{device.status}}).</p>
      </div>
      <div class="panel-footer text-center" v-else>
        <button type="button" class="btn btn-success" title="Reboot this device now."
                v-on:click="rebootDevice()">
          Reboot
        </button>
        <button type="button" class="btn btn-warning"
                v-bind:class="{ disabled: !device.has_config }"
                v-bind:title="device.has_config ? 'Upload configuration and reboot device.' : 'No configuration for this device available.'"
                v-on:click="provisionDevice()">
          Provision
        </button>
        <button type="button" class="btn btn-danger"
                v-bind:class="{ disabled: !device.can_upgrade }"
                v-bind:title="device.can_upgrade ? 'Upload new firmware image and reboot device.' : 'No firmware upgrade for this device available.'"
                v-on:click="upgradeDevice()">
          Firmware Upgrade
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import filters from "../utils/filters.js"

export default {
  name: "DeviceView",
  props: {
    device: {
      type: Object,
      requires: true
    }
  },
  filters: filters,
  methods: {
    rebootDevice: function() {
      this.$emit("device-action", "reboot", this.device.mac_address)
    },
    provisionDevice: function() {
      this.$emit("device-action", "provision", this.device.mac_address)
    },
    upgradeDevice: function() {
      this.$emit("device-action", "upgrade", this.device.mac_address)
    },
    closeView: function() {
      this.$emit("close-view")
    }
  }
}
</script>

<style scoped>

</style>
