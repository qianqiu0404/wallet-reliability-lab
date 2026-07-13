import { describe, expect, it } from 'vitest'
import { advanceSnapshot, createInitialSnapshot } from './useLabEngine'

describe('deterministic lab engine', () => {
  it('starts before the first event', () => {
    expect(createInitialSnapshot('normal-withdrawal')).toEqual({
      scenarioId: 'normal-withdrawal',
      currentIndex: -1,
      playing: false,
      speed: 1,
      mode: 'simulation',
    })
  })

  it('advances in order and stops on the final event', () => {
    let snapshot = { ...createInitialSnapshot('normal-withdrawal'), playing: true }
    snapshot = advanceSnapshot(snapshot, 2)
    expect(snapshot.currentIndex).toBe(0)
    expect(snapshot.playing).toBe(true)
    snapshot = advanceSnapshot(snapshot, 2)
    expect(snapshot.currentIndex).toBe(1)
    expect(snapshot.playing).toBe(false)
  })

  it('does not move beyond the final event', () => {
    const snapshot = { ...createInitialSnapshot('normal-withdrawal'), currentIndex: 1, playing: true }
    expect(advanceSnapshot(snapshot, 2)).toEqual({ ...snapshot, playing: false })
  })
})

