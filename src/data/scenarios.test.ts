import { describe, expect, it } from 'vitest'
import { evidenceRecords, scenarios } from './scenarios'

describe('scenario truth model', () => {
  it('contains exactly the three first-release scenarios', () => {
    expect(scenarios.map(item => item.id)).toEqual([
      'normal-withdrawal',
      'retryable-broadcast',
      'compensation-recovery',
    ])
  })

  it('keeps retryable broadcast failures signed with no balance mutation', () => {
    const scenario = scenarios.find(item => item.id === 'retryable-broadcast')!
    const failure = scenario.steps.find(item => item.id === 'retry-failure')!
    const preserve = scenario.steps.find(item => item.id === 'retry-preserve')!
    expect(failure.state).toBe('signed')
    expect(preserve.state).toBe('signed')
    expect(preserve.details).toContainEqual({ label: 'balance_mutation', value: 'none' })
  })

  it('uses the exact compensation threshold boundary', () => {
    const scenario = scenarios.find(item => item.id === 'compensation-recovery')!
    const fresh = scenario.steps.find(item => item.id === 'comp-wait')!
    const stale = scenario.steps.find(item => item.id === 'comp-stale')!
    expect(fresh.details).toContainEqual({ label: 'eligible', value: 'false' })
    expect(stale.details).toContainEqual({ label: 'eligible', value: 'true' })
  })

  it('resolves every evidence reference', () => {
    const ids = new Set(evidenceRecords.map(item => item.id))
    for (const scenario of scenarios) {
      for (const step of scenario.steps) {
        for (const evidenceId of step.evidenceIds) expect(ids.has(evidenceId)).toBe(true)
      }
    }
  })

  it('links every evidence record to public source evidence', () => {
    for (const evidence of evidenceRecords) {
      expect(evidence.sourceUrl).toMatch(/^https:\/\/github\.com\/qianqiu0404\//)
      expect(evidence.sourceLabel.length).toBeGreaterThan(0)
    }
  })
})
