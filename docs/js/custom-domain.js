document$.subscribe(function () {
  var usernameElem = document.getElementById("username")
  var aliasElem = document.getElementById("alias")

  if (!usernameElem || !aliasElem) {
    return
  }

  var refreshTemplates = function () {
    var q = new URLSearchParams(location.search)
    var username = q.get("username")
    var alias = q.get("alias")

    document.querySelectorAll("code").forEach(function (elem) {
      var template = elem.getAttribute("template")
      if (!template) {
        template = elem.innerHTML
        elem.setAttribute("template", template)
      }


      if (username) {
        var userChunks = username.split("@")
        template = template.replaceAll("__USER_NAME__", userChunks[0])
        template = template.replaceAll("__USER_DOMAIN__", userChunks[1])
      }

      if (alias) {
        var aliasChunks = alias.split("@")
        template = template.replaceAll("__ALIAS_NAME__", aliasChunks[0])
        template = template.replaceAll("__ALIAS_DOMAIN__", aliasChunks[1])
      }

      elem.innerHTML = template
    })
  }

  function isValidUsername(val) {
    if (val.indexOf("@") < 0) {
      return false
    }

    var chunks = val.split("@")
    return isValidHostname(chunks[1])
  }

  var onChange = function () {
    var usernameValue = usernameElem.value
    var aliasValue = aliasElem.value
    var params = new URLSearchParams(location.search)
    var updateURL = false

    if (isValidUsername(usernameValue)) {
      updateURL = true
      params.set("username", usernameValue)

      usernameElem.value = usernameValue
      usernameElem.classList.remove("bad")
    } else {
      usernameElem.classList.add("bad")
    }

    if (isValidUsername(aliasValue)) {
      updateURL = true
      params.set("alias", aliasValue)

      aliasElem.value = aliasValue
      aliasElem.classList.remove("bad")
    } else {
      aliasElem.classList.add("bad")
    }

    if (updateURL) {
      var newUrl = window.location.origin + window.location.pathname + "?" + params.toString()
      window.history.pushState({ path: newUrl }, "", newUrl)
    }

    refreshTemplates()
  }

  var q = new URLSearchParams(location.search)
  usernameElem.onchange = onChange
  usernameElem.value = q.get("username")

  aliasElem.onchange = onChange
  aliasElem.value = q.get("alias")

  onChange()
})
