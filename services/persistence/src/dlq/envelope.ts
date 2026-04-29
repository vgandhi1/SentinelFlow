/**
 * Envelope for persistence DLQ: store reason, timestamp, and original payload as JSON body.
 */
export interface PersistenceDLQEnvelope {
  reason: string;
  error_code?: string;
  timestamp: string; // ISO
  raw_payload: unknown;
}

export function buildDLQEnvelope(reason: string, rawPayload: unknown, errorCode?: string): PersistenceDLQEnvelope {
  return {
    reason,
    error_code: errorCode,
    timestamp: new Date().toISOString(),
    raw_payload: rawPayload,
  };
}
