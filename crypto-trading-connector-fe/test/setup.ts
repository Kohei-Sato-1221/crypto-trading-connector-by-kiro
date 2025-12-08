import { setup, $fetch } from '@nuxt/test-utils/e2e'
import { beforeAll } from 'vitest'

beforeAll(async () => {
  await setup({
    rootDir: process.cwd(),
    server: true,
  })
})
