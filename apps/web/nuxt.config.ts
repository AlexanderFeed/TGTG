export default defineNuxtConfig({
  compatibilityDate: '2026-08-01',
  devtools: { enabled: true },
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
    public: {
      appName: 'ЕщёЕсть',
    },
  },
})
