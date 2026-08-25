// This is the bridge between browser-side Vue code and the Go API.
//
// Nuxt creates a server route from this filename:
//   browser POST /api/v1/auth/login/request
//   -> this catch-all handler (path = "v1/auth/login/request")
//   -> Go POST http://api:8080/v1/auth/login/request
//
// Keeping the request same-origin (/api on the same host as the page) means the
// browser can use one HttpOnly session cookie without exposing the private Go
// container address or needing frontend CORS configuration.
export default defineEventHandler((event) => {
  // apiInternalUrl is private runtime config from nuxt.config.ts. It is read on
  // the Nuxt server and is never included in browser JavaScript.
  const config = useRuntimeConfig(event)
  const path = getRouterParam(event, 'path') || ''

  // Preserve filters such as ?q=bread&limit=20 when forwarding offer requests.
  const query = getRequestURL(event).search
  const base = config.apiInternalUrl.replace(/\/$/, '')

  // proxyRequest forwards method, JSON body, headers, cookies, status, and
  // Set-Cookie. This is why the Go session cookie reaches the browser unchanged.
  return proxyRequest(event, `${base}/${path}${query}`)
})
