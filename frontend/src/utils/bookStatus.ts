// src/utils/bookStatus.ts

export const BOOK_STATUS_LABELS: Record<string, string> = {
    available: 'Available',
    checked_out: 'Checked Out',
}

export const BOOK_STATUS_VARIANTS: Record<string, 'success' | 'error' | 'neutral'> = {
    available: 'success',
    checked_out: 'error',
}

export function getBookStatusLabel(status: string | undefined): string {
    if (!status) return ''
    const key = status.toLowerCase().trim()
    return BOOK_STATUS_LABELS[key] || key.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase())
}

export function getBookStatusVariant(status: string): 'success' | 'error' | 'neutral' {
    if (!status) return 'neutral'
    const key = status?.toLowerCase().trim()
    return BOOK_STATUS_VARIANTS[key] || 'neutral'
}
