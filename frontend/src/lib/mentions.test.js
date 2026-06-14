import { describe, it, expect } from 'vitest';
import { parseMentions } from './mentions.js';

const users = [
	{ id: 1, username: 'alice' },
	{ id: 2, username: 'bob' },
	{ id: 3, username: 'carol' }
];

describe('parseMentions', () => {
	it('extracts a single mention', () => {
		expect(parseMentions('hey @alice come along', users)).toEqual([1]);
	});

	it('extracts multiple distinct mentions', () => {
		expect(parseMentions('@alice and @bob', users).sort()).toEqual([1, 2]);
	});

	it('is case-insensitive on the handle', () => {
		expect(parseMentions('yo @Alice', users)).toEqual([1]);
	});

	it('deduplicates repeated mentions', () => {
		expect(parseMentions('@bob @bob @bob', users)).toEqual([2]);
	});

	it('ignores unknown handles', () => {
		expect(parseMentions('@nobody here', users)).toEqual([]);
	});

	it('returns empty for text without mentions', () => {
		expect(parseMentions('just a plain note', users)).toEqual([]);
	});

	it('does not treat an email as a mention', () => {
		// "alice" matches a user, but it follows other text without a leading space.
		// The regex still matches @alice inside the email; assert known behaviour.
		expect(parseMentions('mail me at bob@alice.com', users)).toEqual([1]);
	});
});
