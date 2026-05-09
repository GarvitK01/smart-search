export async function searchDocuments(query) {
  const res = await fetch(`/search?q=${encodeURIComponent(query)}`);
  if (!res.ok) throw new Error("Search failed");
  return res.json();
}

export async function uploadFile(file) {
  const form = new FormData();
  form.append("file", file);
  const res = await fetch("/upload", { method: "POST", body: form });
  if (!res.ok) throw new Error("Upload failed");
  return res.text();
}

export async function listDocuments() {
  const res = await fetch("/documents");
  if (!res.ok) throw new Error("Failed to load documents");
  return res.json();
}

export async function login(username, password) {
  const res = await fetch("/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  });
  if (!res.ok) {
    const msg = await res.text();
    throw new Error(msg || "Login failed");
  }
  return res.text();
}

export async function register(username, password) {
  const res = await fetch("/register", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  });
  if (!res.ok) {
    const msg = await res.text();
    throw new Error(msg || "Registration failed");
  }
  return res.text();
}

export async function logout() {
  const res = await fetch("/logout", { method: "POST" });
  if (!res.ok) throw new Error("Logout failed");
  return res.text();
}

export async function checkAuth() {
  const res = await fetch("/me");
  if (!res.ok) throw new Error("Not authenticated");
  return res.json();
}
