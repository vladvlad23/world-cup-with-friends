/** Extract user ids whose @username appears in the text. */
export function parseMentions(text, users) {
	const found = new Set();
	const matches = text.matchAll(/@([a-zA-Z0-9_]+)/g);
	for (const m of matches) {
		const handle = m[1].toLowerCase();
		const user = users.find((u) => u.username.toLowerCase() === handle);
		if (user) found.add(user.id);
	}
	return [...found];
}
