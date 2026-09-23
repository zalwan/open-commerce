// Shared UI helpers (no API calls here).
export function statusClass(s: string): string {
	if (s === 'paid' || s === 'done') return 'badge badge-ok';
	if (s === 'pending' || s === 'shipped') return 'badge badge-info';
	if (s === 'payment_failed') return 'badge badge-err';
	return 'badge';
}
