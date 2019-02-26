function createNameElement(name) {
  var div = document.createElement("div");
  div.innerHTML = `Name: ${name}`;
  div.className = "endpointEl-name";

  return div;
}

function createEndpointValueElement(endpoint) {
  var div = document.createElement("div");
  div.innerHTML = `Endpoint: ${endpoint}`;
  div.className = "endpointEl-endpoint-value";

  return div;
}

function createStatusElement(lastPing) {
  // Calcluating time since last ping
  var diff = new Date() - Date.parse(lastPing);
  var status = diff < 10000 ? `Status: Healthy` : `Status: Unhealthy`;

  var div = document.createElement("div");
  div.innerHTML = status;
  div.className = "endpointEl-status";

  return div;
}

function createLastPingedElement(lastPing) {
  var date = new Date(Date.parse(lastPing));
  var div = document.createElement("div");
  if (date.getFullYear() === 1) {
    div.innerHTML = `Last Successful Ping: Never`;
  } else {
    var options = {
      weekday: "long",
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit"
    };
    div.innerHTML = `Last Successful Ping: ${date.toLocaleString(
      "en-US",
      options
    )}`;
  }

  div.className = "endpointEl-last-ping";

  return div;
}

function createDeletButton(endpoint) {
  var btn = document.createElement("BUTTON");
  var t = document.createTextNode("Delete Site");
  btn.onclick = removeSite;
  btn.id = endpoint;
  btn.appendChild(t);
  return btn;
}

function removeSite(e) {
  e.preventDefault();

  var xhttp = new XMLHttpRequest();

  xhttp.onreadystatechange = function() {
    if (this.readyState === 4) {
      if (this.status === 202) {
        var data = JSON.parse(this.responseText);
        var endpoints = Object.keys(data);

        if (endpoints.length === 0) {
          var sitesEl = document.getElementById("sites");
          sitesEl.innerHTML = "";
        } else if (endpoints.length > 0) {
          var sitesEl = document.getElementById("sites");
          sitesEl.innerHTML = "";

          for (var i = 0; i < endpoints.length; i++) {
            var endpointEl = createEndpointElement({ ...data[endpoints[i]] });
            sitesEl.appendChild(endpointEl);
          }
        }
      } else {
        var data = JSON.parse(this.responseText);
        errorsElement = document.getElementById("errors");
        errorsElement.style.display = "block";
        if (data.errors) {
          errorsElement.innerHTML = data.errors.join(", ");
        } else {
          errorsElement.innerHTML = data.error;
        }
      }
    }
  };

  var data = { endpoint: e.target.id };

  xhttp.open("DELETE", "/api/sites", true);
  xhttp.send(JSON.stringify(data));
}

function createEndpointElement(data) {
  var endpointEl = document.createElement("div");
  endpointEl.className = "endpointEl-container";

  var nameEl = createNameElement(data.name);
  endpointEl.appendChild(nameEl);

  var endpointValEl = createEndpointValueElement(data.endpoint);
  endpointEl.appendChild(endpointValEl);

  var statusEl = createStatusElement(data.lastPing);
  endpointEl.appendChild(statusEl);

  var lastPingEl = createLastPingedElement(data.lastPing);
  endpointEl.appendChild(lastPingEl);

  var deleteButton = createDeletButton(data.endpoint);
  endpointEl.appendChild(deleteButton);

  return endpointEl;
}

function addSite(e) {
  e.preventDefault();

  var xhttp = new XMLHttpRequest();

  xhttp.onreadystatechange = function() {
    if (this.readyState === 4) {
      if (this.status === 200) {
        errorsElement = document.getElementById("errors");
        errorsElement.style.display = "none";
        var sitesEl = document.getElementById("sites");
        var data = JSON.parse(this.responseText);
        var endpointEl = createEndpointElement(data);
        sitesEl.appendChild(endpointEl);
      } else {
        var data = JSON.parse(this.responseText);
        errorsElement = document.getElementById("errors");
        errorsElement.style.display = "block";
        if (data.errors) {
          errorsElement.innerHTML = data.errors.join(", ");
        } else {
          errorsElement.innerHTML = data.error;
        }
      }
    }
  };

  var name = document.getElementById("name").value;
  var endpoint = document.getElementById("endpoint").value;
  var data = { name: name, endpoint: endpoint };

  xhttp.open("POST", "/api/sites", true);
  xhttp.send(JSON.stringify(data));
}

function getSites() {
  var xhttp = new XMLHttpRequest();

  xhttp.onreadystatechange = function() {
    if (this.readyState === 4) {
      if (this.status === 200) {
        var data = JSON.parse(this.responseText);
        var endpoints = Object.keys(data);

        if (endpoints.length > 0) {
          var sitesEl = document.getElementById("sites");

          for (var i = 0; i < endpoints.length; i++) {
            var endpointEl = createEndpointElement({ ...data[endpoints[i]] });
            sitesEl.appendChild(endpointEl);
          }
        }
      }
    }
  };

  xhttp.open("GET", "/api/sites", true);
  xhttp.send();
}
