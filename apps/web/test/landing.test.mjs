import assert from 'node:assert/strict'
import { readFile } from 'node:fs/promises'
import test from 'node:test'

const page = await readFile(new URL('../src/pages/index.astro', import.meta.url), 'utf8')

test('landing explains the secure env-to-GitHub workflow', () => {
  assert.match(page, /Preview every change before it reaches GitHub/)
  assert.match(page, /env\.production\.to\.github/)
  assert.match(page, /etg --dry-run/)
})

test('landing exposes clear installation and GitHub calls to action', () => {
  assert.match(page, /href="#install"/)
  assert.match(page, /https:\/\/github\.com\/mattrbdr\/envtogithub/)
})

test('landing uses a green-only visual system and component icons', () => {
  assert.match(page, /#c8f135/)
  assert.match(page, /from '@lucide\/astro'/)
  assert.match(page, /aria-label="GitHub"/)
  assert.doesNotMatch(page, /linear-gradient|radial-gradient/)
  assert.doesNotMatch(page, /#f6c945/)
})

test('workflow keeps its copy in a full-width step column', () => {
  assert.match(page, /\.step \{ border-left:1px solid #435136; display:block;/)
})

test('installation points to published release archives instead of an unavailable tap', () => {
  assert.match(page, /releases\/download\/v0\.1\.3/)
  assert.match(page, /etg_0\.1\.3_darwin_arm64\.tar\.gz/)
  assert.match(page, /Install from the latest release/)
  assert.doesNotMatch(page, /brew tap mattrbdr\/tap/)
})
