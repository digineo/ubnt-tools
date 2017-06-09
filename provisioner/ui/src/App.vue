<template>
  <div class="container-fluid" id="app">
    <ubnt-provisioner
        v-bind:provisioner="provisioner">
    </ubnt-provisioner>
  </div>
</template>

<script>
import UbntProvisioner  from "./components/UbntProvisioner.vue"

export default {
  name: "App",
  props: {
    provisioner: {
      type: Object,
      required: true
    }
  },
  components: {
    UbntProvisioner
  }
}
</script>
