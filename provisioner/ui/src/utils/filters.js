import moment from "moment"

let filters = {
  timeAgo(value) {
    let unix = moment.unix(value),
        diff = moment().unix() - unix.unix()
    return diff < 60 ? `${diff} seconds ago` : unix.fromNow()
  },

  fmtDate(value) {
    let unix = moment.unix(value),
        fmt = unix.format("YYYY-MM-DD, HH:mm"),
        diff = moment().unix() - unix.unix(),
        ago = diff < 60 ? `${diff} seconds ago` : unix.fromNow()
    return `${fmt} (${ago})`
  }
}

export default filters
