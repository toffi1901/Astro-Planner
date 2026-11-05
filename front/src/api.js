const BASE = 'http://localhost:8080/api'
let accessToken = null

export const setAccessToken = (t) => { accessToken = t}

async function query(path, options = {}) {
  const headers = { 'Content-Type': 'application/json', ...options.headers }
  if (accessToken) headers.Authorization = `Bearer ${accessToken}`
  return fetch(BASE + path, { ...options, headers, credentials: 'include'})
}

export async function refresh() {
  const r = await query('/auth/refresh', { method: 'POST' })
  if (!r.ok) {
    setAccessToken(null)
    return false
  }
  setAccessToken((await r.json()).access_token)
  return true
}

export async function login(email, password) {
  const r = await query('/auth/login', {
    method: 'POST',
    body: JSON.stringify({email, password})
  })
  const data = await r.json()
  if (!r.ok) throw new Error(data.error)
  setAccessToken(data.access_token)
  return data
}

export async function reg(email, password) {
  const r = await query('/auth/register', {
    method: 'POST',
    body: JSON.stringify({email, password})
  })
  const data = await r.json()
  if (!r.ok) throw new Error(data.error)
  return data
}

export async function apiFetch(path, options) {
  let r = await query(path, options)
  if (r.status == 401 && await refresh()) r = await query(path, options)
  return r
}