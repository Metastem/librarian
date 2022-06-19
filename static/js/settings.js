function getCookie(name) {
  let cookies = document.cookie.split("; ");
  let cookie = cookies.filter(cookie => cookie.includes(name))[0]
  return cookie ? cookie.split("=")[1] : "";
}

// General settings
let generalSettings = ["theme", "showRelated", "collapseComments", "nsfw"]

function updateGeneralSetting(setting) {
  let elem = document.getElementById(setting)

  if (setting == "nsfw") {
    document.cookie = "nsfw=" + elem.checked + "; path=/; SameSite=Strict; max-age=2147483647"
  } else if (setting == "theme") {
    document.cookie = "theme=" + elem.value + "; path=/; SameSite=Strict; max-age=2147483647"
    location.reload()
  } else {
    localStorage.setItem(setting, elem.checked)
  }
}

function loadGeneralSetting(setting) {
  let elem = document.getElementById(setting)

  if (setting == "nsfw") {
    if (document.cookie.includes("nsfw=true")) {
      elem.checked = true
    }
  } else if (setting == "theme") {
    elem.value = getCookie("theme") || "system"
  } else {
    let value = localStorage.getItem(setting)
    if (value) {
      elem.checked = value
    }
  }
}

generalSettings.forEach(setting => {
  let select = document.getElementById(setting)
  select.addEventListener("change", () => updateGeneralSetting(setting))
  loadGeneralSetting(setting)
})

// Player settings
let autoplaySelect = document.getElementById("autoplay");
let autoplayNextVidSelect = document.getElementById("autoplayNextVid");
let speedSelect = document.getElementById("speed");
let qualitySelect = document.getElementById("quality");

function loadSettings() {
  let plyrSettings = JSON.parse(localStorage.getItem("plyr")) || {};
  speedSelect.value = plyrSettings.speed || "1";
  qualitySelect.value = plyrSettings.quality || "0";
  autoplaySelect.checked = localStorage.getItem("autoplay") === "true" || false;
  autoplayNextVidSelect.checked = localStorage.getItem("autoplayNextVid") === "true" || false;
}
loadSettings()

autoplaySelect.addEventListener("change", () => {
  localStorage.setItem("autoplay", autoplaySelect.checked);
})
autoplayNextVidSelect.addEventListener("change", () => {
  localStorage.setItem("autoplayNextVid", autoplayNextVidSelect.checked);
})
qualitySelect.addEventListener("change", () => {
  let plyrStorage = localStorage.getItem("plyr");
  let data = plyrStorage ? JSON.parse(plyrStorage) : {};
  data.quality = qualitySelect.value * 1;
  localStorage.setItem("plyr", JSON.stringify(data));
})
speedSelect.addEventListener("change", () => {
  let plyrStorage = localStorage.getItem("plyr");
  let data = plyrStorage ? JSON.parse(plyrStorage) : {};
  data.speed = speedSelect.value * 1;
  localStorage.setItem("plyr", JSON.stringify(data));
})

// SponsorBlock settings
let categories = ["sponsor", "selfpromo", "interaction", "intro", "outro", "preview", "filler"]

function updateSBSetting(category) {
  let categories = localStorage.getItem("sb_categories") || "";
  if (categories.includes(category)) {
    let re = new RegExp(`,?${category}`)
    localStorage.setItem("sb_categories", categories.replace(re, ""));
  } else if (categories.length == 0) {
    localStorage.setItem("sb_categories", categories + category);
  } else {
    localStorage.setItem("sb_categories", categories + "," + category);
  }

  let newCategories = localStorage.getItem("sb_categories")
  if (newCategories.startsWith(",")) {
    localStorage.setItem("sb_categories", newCategories.substring(1, 999));
  }
}

categories.forEach(category => {
  let select = document.getElementById(category)
  select.addEventListener("change", () => updateSBSetting(category))
  let categories = localStorage.getItem("sb_categories") || ""
  if (categories.includes(category)) {
    select.checked = true
  }
})