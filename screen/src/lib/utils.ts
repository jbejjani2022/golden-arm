// src/lib/utils.ts
export function formatDate(dateString: string): string {
    const date = new Date(dateString);
    const options: Intl.DateTimeFormatOptions = {
        month: '2-digit',
        day: '2-digit',
        year: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        hour12: true,
    };
    return date.toLocaleString('en-US', options);
}

export function formatDateFriendly(dateString: string): string {
    const date = new Date(dateString);
    const weekday = date.toLocaleDateString('en-US', { weekday: 'long' });
    const month = date.toLocaleDateString('en-US', { month: 'short' });
    const day = date.getDate();
    const year = date.getFullYear();
    const time = date.toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit', hour12: true });
    return `${weekday}, ${month}. ${day}, ${year}, ${time}`;
}

export function formatRuntime(runtime: number): string {
    if (runtime < 0) {
        throw new Error("Runtime cannot be negative");
    }

    const hours = Math.floor(runtime / 60);
    const minutes = runtime % 60;

    if (hours > 0) {
        return `${hours}h ${minutes}m`;
    } else {
        return `${minutes}m`;
    }
}
