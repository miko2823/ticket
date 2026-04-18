import { useEffect, useRef, useState } from "react";
import { api } from "../api/client";

interface SessionResponse {
  session_id: string;
  heartbeat_interval_ms: number;
}

export function useSession(eventId: string) {
  const sessionRef = useRef<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isConnecting, setIsConnecting] = useState(true);

  useEffect(() => {
    let heartbeatTimer: number;
    let cancelled = false;

    const startSession = async () => {
      try {
        const res = await api.post<SessionResponse>(
          `/events/${eventId}/session`,
          {}
        );
        if (cancelled) return;
        sessionRef.current = res.session_id;
        setIsConnecting(false);

        heartbeatTimer = window.setInterval(async () => {
          if (!sessionRef.current) return;
          try {
            await api.put(
              `/events/${eventId}/session/${sessionRef.current}`,
              {}
            );
          } catch {
            // Session expired — stop heartbeat
            clearInterval(heartbeatTimer);
            sessionRef.current = null;
            setError("Session expired. Please reopen the seat map.");
          }
        }, res.heartbeat_interval_ms);
      } catch (err) {
        if (cancelled) return;
        setIsConnecting(false);
        setError(
          err instanceof Error ? err.message : "Failed to create session"
        );
      }
    };

    startSession();

    const handleUnload = () => {
      // Best-effort cleanup — TTL handles the rest
      if (sessionRef.current) {
        api
          .del(`/events/${eventId}/session/${sessionRef.current}`)
          .catch(() => {});
      }
    };
    window.addEventListener("beforeunload", handleUnload);

    return () => {
      cancelled = true;
      clearInterval(heartbeatTimer);
      window.removeEventListener("beforeunload", handleUnload);
      if (sessionRef.current) {
        api
          .del(`/events/${eventId}/session/${sessionRef.current}`)
          .catch(() => {});
        sessionRef.current = null;
      }
    };
  }, [eventId]);

  return { sessionId: sessionRef, error, isConnecting };
}
