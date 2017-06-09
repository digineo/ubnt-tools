import Vue          from "vue"
import jQuery       from "jquery"
import Provisioner  from "./provisioner.js"
import App          from "./App.vue"

const apiEndpoint = "/api"

let directory = jQuery.getJSON(apiEndpoint)

directory.done((data, _status, _xhr)=>{
  let prov = new Provisioner(data)
  new Vue({
    el: "#app",
    template: `<App :provisioner="provisioner" />`,
    data: {
      provisioner: prov,
    },
    components: {
      App
    }
  })
})

directory.fail(()=>{
  jQuery("#app").text("Failed to initialize Provisioner.")
})
