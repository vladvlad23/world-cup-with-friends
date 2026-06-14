import { describe, it, expect } from 'vitest';
import { dayKey, timeLabel, prettyDate, buildMonth, monthName } from './date.js';

describe('dayKey / timeLabel (Bucharest time)', () => {
	it('converts an instant to the Bucharest day', () => {
		// 13:00Z in June -> 16:00 EEST, same day.
		expect(dayKey('2026-06-11T13:00:00Z')).toBe('2026-06-11');
	});
	it('converts an instant to the Bucharest clock time', () => {
		expect(timeLabel('2026-06-11T13:00:00Z')).toBe('16:00');
	});
	it('rolls a late-night US kickoff onto the next Bucharest day', () => {
		// 22:00Z -> 01:00 next day in Bucharest (EEST, UTC+3).
		expect(dayKey('2026-06-23T22:00:00Z')).toBe('2026-06-24');
		expect(timeLabel('2026-06-23T22:00:00Z')).toBe('01:00');
	});
});

describe('prettyDate', () => {
	it('formats a day key with weekday', () => {
		expect(prettyDate('2026-06-11')).toBe('Thu, June 11, 2026');
	});
});

describe('monthName', () => {
	it('maps a zero-based index', () => {
		expect(monthName(5)).toBe('June');
		expect(monthName(6)).toBe('July');
	});
});

describe('buildMonth', () => {
	it('lays June 2026 out on the correct weekdays', () => {
		const weeks = buildMonth(2026, 5);
		// June 1, 2026 is a Monday -> one leading blank (Sunday).
		expect(weeks[0][0]).toBeNull();
		expect(weeks[0][1]).toEqual({ day: 1, key: '2026-06-01' });
	});

	it('covers every day of the month exactly once', () => {
		const weeks = buildMonth(2026, 5); // June has 30 days
		const days = weeks.flat().filter(Boolean).map((c) => c.day);
		expect(days).toEqual(Array.from({ length: 30 }, (_, i) => i + 1));
	});

	it('pads each week to 7 cells', () => {
		for (const week of buildMonth(2026, 6)) {
			expect(week).toHaveLength(7);
		}
	});
});
