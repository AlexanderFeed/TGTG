export default defineEventHandler((event) => {
  const config = useRuntimeConfig(event)
  const path = getRouterParam(event, 'path') || ''
  const query = getRequestURL(event).search
  const base = config.apiInternalUrl.replace(/\/$/, '')

  return proxyRequest(event, `${base}/${path}${query}`)
})
