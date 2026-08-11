import assert from 'node:assert/strict'
import { readFile } from 'node:fs/promises'
import test from 'node:test'

const workflow = await readFile(new URL('../../../.github/workflows/deploy-on-production.yaml', import.meta.url), 'utf8')

test('production backup has bounded SSH execution and only enables rollback after it succeeds', () => {
  assert.match(workflow, /timeout 180s ssh/)
  assert.match(workflow, /ConnectTimeout=20/)
  assert.match(workflow, /echo "name=\$BACKUP_NAME" >> "\$GITHUB_OUTPUT"/)
  assert.doesNotMatch(workflow, /BACKUP_NAME=.*\n\s*echo "name=\$BACKUP_NAME"[\s\S]*?timeout 180s ssh/)
})

test('production deployment verifies SSH access after whitelisting before backing up', () => {
  assert.match(workflow, /name: Verify SSH connectivity/)
  assert.match(workflow, /for attempt in \{1\.\.6\}/)
  assert.match(workflow, /SSH remained unreachable after whitelisting/)
})
