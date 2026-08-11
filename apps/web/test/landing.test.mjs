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

test('installation uses the published Homebrew tap', () => {
  assert.match(page, /brew tap mattrbdr\/tap/)
  assert.match(page, /brew install etg/)
  assert.match(page, /Install with Homebrew/)
})
