// Thunderbird can terminate idle backgrounds in manifest v3.
// Any listener directly added during add-on startup will be registered as a
// persistent listener and the background will wake up (restart) each time the
// event is fired.

browser.composeAction.onClicked.addListener(async (tab) => {
  // Get the existing message.
  let details = await browser.compose.getComposeDetails(tab.id);

  let params = new URLSearchParams();

  if (!details.to || !details.from || !details.subject) {
    let document = new DOMParser().parseFromString(details.body, "text/html");
    let para = document.createElement("p");
    para.textContent =
      "Please fill in all fields before adding tracking information.";
    document.body.appendChild(para);
    let html = new XMLSerializer().serializeToString(document);
    browser.compose.setComposeDetails(tab.id, { body: html });

    return;
  }

  params.append("to", details.to);
  params.append("from", details.from);
  params.append("subject", details.subject);

  let token = await fetch(`https://t.jackmerrill.com/t?${params.toString()}`, {
    method: "GET",
    mode: "cors",
  }).then((r) => r.text());

  // The message is being composed in HTML mode. Parse the message into an HTML document.
  let document = new DOMParser().parseFromString(details.body, "text/html");
  console.log(document);

  // Use normal DOM manipulation to modify the message.
  let para = document.createElement("img");
  para.src = `https://t.jackmerrill.com/i?id=${token}`;
  para.width = 1;
  para.height = 1;
  document.body.appendChild(para);

  // Serialize the document back to HTML, and send it back to the editor.
  let html = new XMLSerializer().serializeToString(document);
  console.log(html);
  browser.compose.setComposeDetails(tab.id, { body: html });
});
