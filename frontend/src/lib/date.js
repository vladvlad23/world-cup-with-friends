// The backend stores each kickoff as a correct absolute UTC instant. We display
// everything in Bucharest time (all users are in Romania), converting both the
// calendar day and the clock time so a late US night match lands on the right day.
export const DISPLAY_TZ = 'Europe/Bucharest';

const dayFmt = new Intl.DateTimeFormat('en-CA', {
	timeZone: DISPLAY_TZ,
	year: 'numeric',
	month: '2-digit',
	day: '2-digit'
});
const timeFmt = new Intl.DateTimeFormat('en-GB', {
	timeZone: DISPLAY_TZ,
	hour: '2-digit',
	minute: '2-digit',
	hour12: false
});

/** ISO instant -> "YYYY-MM-DD" in Bucharest time (en-CA formats as YYYY-MM-DD). */
export function dayKey(iso) {
	return dayFmt.format(new Date(iso));
}

/** ISO instant -> "HH:MM" in Bucharest time. */
export function timeLabel(iso) {
	return timeFmt.format(new Date(iso));
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
