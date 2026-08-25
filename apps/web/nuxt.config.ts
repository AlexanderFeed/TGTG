export default defineNuxtConfig({
  compatibilityDate: '2026-08-01',
  devtools: { enabled: false },
  css: ['~/assets/css/main.css'],
  typescript: {
    strict: true,
    typeCheck: true,
  },
  app: {
    head: {
      htmlAttrs: { lang: 'ru' },
      title: 'ЕщёЕсть — спасаем хорошую еду',
      meta: [
        {
          name: 'description',
          content: 'Забирайте хорошую еду из любимых мест со скидкой и помогайте ей не пропасть.',
        },
        { name: 'theme-color', content: '#173e32' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1, viewport-fit=cover' },
      ],
      link: [{ rel: 'icon', type: 'image/svg+xml', href: '/favicon.svg' }],
    },
  },
  runtimeConfig: {
    // Private server-only value used by server/api/[...path].ts. In Docker it can
    // be overridden with NUXT_API_INTERNAL_URL=http://api:8080. Do not move it to
    // `public`, because browsers cannot resolve the Docker service name `api`.
    apiInternalUrl: 'http://127.0.0.1:8080',
    public: {
      // Values under public are intentionally serialized into browser JavaScript.
      appName: 'ЕщёЕсть',
    },
  },
})
