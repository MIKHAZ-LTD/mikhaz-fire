// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2026-07-02',
  ssr: false,
  modules: ['@nuxt/ui'],
  css: ['~/assets/css/main.css'],
  router: {
    options: {
      hashMode: true
    }
  },
  nitro: {
    output: {
      publicDir: 'dist'
    }
  }
})
