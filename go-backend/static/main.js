function getCookie(cname) {
  let name = cname + "=";
  let ca = document.cookie.split(';');

  for(let i = 0; i < ca.length; i++) {
    let c = ca[i];
    
    while (c.charAt(0) == ' ') {
      c = c.substring(1);
    }

    if (c.indexOf(name) == 0) {
      return c.substring(name.length, c.length);
    }

    
  }
  return "NOT FOUND";
}

function checkCookie() {
  let cookievalue = getCookie("cars_viewer_allow_tracking");

  if (cookievalue == "true" || cookievalue == "false") {
    console.log("cookie allowance: " + cookievalue);
    console.log("long cookie: " + getCookie("cars_viewer_long_term_id") + "<-- http only cookie");
    console.log("short cookie: " + getCookie("cars_viewer_short_term_id"));

  } else {
    document.getElementById("cookie_div").style.display = "block";
  }
}

window.addEventListener("DOMContentLoaded", checkCookie);


document.getElementById("accept_cookies_button").addEventListener("click", () => {
    document.getElementById("cookie_div").style.display = "none";
    
    fetch("/allow-cookies");

});

document.getElementById("decline_cookies_button").addEventListener("click", () => {
    document.getElementById("cookie_div").style.display = "none";
    
    fetch("/disallow-cookies");
});
