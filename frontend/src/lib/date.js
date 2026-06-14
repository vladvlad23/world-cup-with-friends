// The backend stores kickoff times as the match's local wall-clock time encoded
// as UTC (e.g. "2026-06-11T13:00:00Z"). To display the intended day/time without
// the browser's timezone shifting it, we read the ISO string components directly.

/** "2026-06-11T13:00:00Z" -> "2026-06-11" */
export function dayKey(iso) {
	return iso.slice(0, 10);
}

/** "2026-06-11T13:00:00Z" -> "13:00" */
export function timeLabel(iso) {
	return iso.slice(11, 16);
}

const MONTHS = [
	'January', 'February', 'March', 'April', 'May', 'June',
	'July', 'August', 'September', 'October', 'November', 'December'
];
const WEEKDAYS = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];

export const weekdayLabels = WEEKDAYS;

export function monthName(monthIndex) {
	return MONTHS[monthIndex];
}

/** "2026-06-11" -> "Thu, June 11, 2026" */
export function prettyDate(key) {
	const [y, m, d] = key.split('-').map(Number);
	const dt = new Date(Date.UTC(y, m - 1, d));
	return `${WEEKDAYS[dt.getUTCDay()]}, ${MONTHS[m - 1]} ${d}, ${y}`;
}

/**
 * Build a calendar grid (array of weeks; each week is 7 cells) for the given
 * year/month. Cells outside the month are marked inMonth:false.
 */
export function buildMonth(year, monthIndex) {
	const first = new Date(Date.UTC(year, monthIndex, 1));
	const startDay = first.getUTCDay();
	const daysInMonth = new Date(Date.UTC(year, monthIndex + 1, 0)).getUTCDate();

	const cells = [];
	// Leading blanks from previous month.
	for (let i = 0; i < startDay; i++) cells.push(null);
	for (let d = 1; d <= daysInMonth; d++) {
		const key = `${year}-${String(monthIndex + 1).padStart(2, '0')}-${String(d).padStart(2, '0')}`;
		cells.push({ day: d, key });
	}
	while (cells.length % 7 !== 0) cells.push(null);

	const weeks = [];
	for (let i = 0; i < cells.length; i += 7) weeks.push(cells.slice(i, i + 7));
	return weeks;
}
