import { useEffect, useRef, useState } from "react";
import { api } from "../api/client";
import { getRecaptchaToken } from "../api/recaptcha";

interface SessionCreatedResponse {
  session_id: string;
  heartbeat_interval_ms: number;
}

interface QueuedResponse {
  status: "queued";
  position: number;
  estimated_wait_seconds: number;
  queue_length: number;
}

interface QueueStatusResponse {
  status: "waiting" | "admitted" | "none";
  position?: number;
  estimated_wait_seconds?: number;
  queue_length?: number;
}

export interface QueueInfo {
  position: number;
  estimatedWait: number;
  queueLength: number;
}

type CreateResponse = SessionCreatedResponse | QueuedResponse;

function isQueued(res: CreateResponse): res is QueuedResponse {
  return "status" in res && res.status === "queued";
}

export function useSession(eventId: string) {
  const sessionRef = useRef<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isConnecting, setIsConnecting] = useState(true);
  const [queueInfo, setQueueInfo] = useState<QueueInfo | null>(null);

  useEffect(() => {
    let heartbeatTimer: number;
    let pollTimer: number;
    let cancelled = false;

    const startHeartbeat = (sessionId: string, intervalMs: number) => {
      heartbeatTimer = window.setInterval(async () => {
        try {
          await api.put(`/events/${eventId}/session/${sessionId}`, {});
        } catch {
          clearInterval(heartbeatTimer);
          sessionRef.current = null;
          setError("Session expired. Please reopen the seat map.");
        }
      }, intervalMs);
    };

    const createSession = async (): Promise<boolean> => {
      const recaptchaToken = await getRecaptchaToken("create_session");
      const res = await api.post<CreateResponse>(
        `/events/${eventId}/session`,
        {},
        recaptchaToken ? { "X-Recaptcha-Token": recaptchaToken } : undefined
      );
      if (cancelled) return false;

      if (isQueued(res)) {
        setQueueInfo({
          position: res.position,
          estimatedWait: res.estimated_wait_seconds,
          queueLength: res.queue_length,
        });
        return false;
      }

      // Session created
      sessionRef.current = res.session_id;
      setQueueInfo(null);
      setIsConnecting(false);
      startHeartbeat(res.session_id, res.heartbeat_interval_ms);
      return true;
    };

    const startPolling = () => {
      pollTimer = window.setInterval(async () => {
        if (cancelled) return;
        try {
          const status = await api.get<QueueStatusResponse>(
            `/events/${eventId}/queue`
          );
          if (cancelled) return;

          if (status.status === "admitted") {
            clearInterval(pollTimer);
            await createSession();
          } else if (status.status === "waiting") {
            setQueueInfo({
              position: status.position ?? 0,
              estimatedWait: status.estimated_wait_seconds ?? 0,
              queueLength: status.queue_length ?? 0,
            });
          } else {
            // "none" — no longer in queue (shouldn't happen normally)
            clearInterval(pollTimer);
            setIsConnecting(false);
            setError("You are no longer in the queue.");
          }
        } catch {
          // Polling error — keep trying
        }
      }, 3000);
    };

    const init = async () => {
      try {
        const created = await createSession();
        if (cancelled) return;
        if (!created) {
          // Queued — start polling
          setIsConnecting(false);
          startPolling();
        }
      } catch (err) {
        if (cancelled) return;
        setIsConnecting(false);
        setError(
          err instanceof Error ? err.message : "Failed to create session"
        );
      }
    };

    init();

    const handleUnload = () => {
      if (sessionRef.current) {
        api
          .del(`/events/${eventId}/session/${sessionRef.current}`)
          .catch(() => {});
      } else if (queueInfo) {
        api.del(`/events/${eventId}/queue`).catch(() => {});
      }
    };
    window.addEventListener("beforeunload", handleUnload);

    return () => {
      cancelled = true;
      clearInterval(heartbeatTimer);
      clearInterval(pollTimer);
      window.removeEventListener("beforeunload", handleUnload);
      if (sessionRef.current) {
        api
          .del(`/events/${eventId}/session/${sessionRef.current}`)
          .catch(() => {});
        sessionRef.current = null;
      }
    };
  }, [eventId]);

  return { sessionId: sessionRef, error, isConnecting, queueInfo };
}
