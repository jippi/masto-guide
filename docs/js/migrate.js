function url_domain(data) {
  if (isValidHostname(data)) {
    return data
  }

  var a = document.createElement("a")
  a.href = data
  return a.hostname
}

document$.subscribe(function () {
  var fromServerElement = document.getElementById("from-server")
  var toServerElement = document.getElementById("to-server")

  if (!fromServerElement || !toServerElement) {
    return
  }

  var refreshLinks = function () {
    var q = new URLSearchParams(location.search)

    var fromServer = q.get("from")
    var toServer = q.get("to")

    document.querySelectorAll("a[identity]").forEach(function (link) {
      var identity = link.getAttribute("identity")
      if (identity == "old") {
        if (fromServer) {
          link.setAttribute("href", "https://" + fromServer + link.getAttribute("path"))
          link.onclick = null

          return
        }

        link.setAttribute("href", "#nope")
        link.onclick = warnMissingConfig("Please configure your old server in step 1")
      }

      if (identity == "new") {
        if (toServer) {
          link.setAttribute("href", "https://" + toServer + link.getAttribute("path"))
          link.onclick = null

          return
        }

        link.setAttribute("href", "#nope")
        link.onclick = warnMissingConfig("Please configure your new server in step 1")
      }
    })

    document.querySelectorAll("span[identity]").forEach(function (span) {
      var identity = span.getAttribute("identity")
      var innerSpan = document.createElement(span.getAttribute("tag"))

      if (identity == "old") {
        if (fromServer) {
          innerSpan.innerText = fromServer
        } else {
          innerSpan.innerText = span.getAttribute("default")
        }
      }

      if (identity == "new") {
        if (toServer) {
          innerSpan.innerText = toServer
        } else {
          innerSpan.innerText = span.getAttribute("default")
        }
      }

      innerSpan.innerText =
        (span.getAttribute("prefix") ?? "") + innerSpan.innerHTML + (span.getAttribute("suffix") ?? "")

      span.innerHTML = innerSpan.outerHTML
    })
  }

  var onChange = function () {
    var fromServerDomain = url_domain(fromServerElement.value)
    var toServerDomain = url_domain(toServerElement.value)
    var updateURL = false
    var q = new URLSearchParams(location.search)

    if (isValidHostname(fromServerDomain)) {
      updateURL = true
      q.set("from", fromServerDomain)

      fromServerElement.value = fromServerDomain
      fromServerElement.classList.remove("bad")
    } else {
      fromServerElement.classList.add("bad")
    }

    if (isValidHostname(toServerDomain)) {
      updateURL = true
      q.set("to", toServerDomain)

      toServerElement.value = toServerDomain
      toServerElement.classList.remove("bad")
    } else {
      toServerElement.classList.add("bad")
    }

    if (updateURL) {
      var newUrl = window.location.origin + window.location.pathname + "?" + q.toString()
      window.history.pushState({ path: newUrl }, "", newUrl)
    }

    refreshLinks()
  }

  var q = new URLSearchParams(location.search)
  fromServerElement.onchange = onChange
  fromServerElement.value = q.get("from")

  toServerElement.onchange = onChange
  toServerElement.value = q.get("to")

  onChange()
})

function isValidHostname(value) {
  if (typeof value !== "string") return false
  if (value == location.hostname) return false

  const validHostnameChars = /^[a-zA-Z0-9-.]{1,253}\.?$/g
  if (!validHostnameChars.test(value)) {
    return false
  }

  if (value.endsWith(".")) {
    value = value.slice(0, value.length - 1)
  }

  if (value.length > 253) {
    return false
  }

  const labels = value.split(".")

  const isValid = labels.every(function (label) {
    const validLabelChars = /^([a-zA-Z0-9-]+)$/g

    const validLabel =
      validLabelChars.test(label) && label.length < 64 && !label.startsWith("-") && !label.endsWith("-")

    return validLabel
  })

  return isValid
}

function warnMissingConfig(text) {
  return function () {
    alert(text)
    return false
  }
}
