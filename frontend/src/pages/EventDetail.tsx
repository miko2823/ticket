import { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";
import { api } from "../api/client";

interface Event {
  id: string;
  name: string;
  venue: string;
  starts_at: string;
  ticketing_starts_at: string;
  ticketing_ends_at: string;
}

function formatCountdown(seconds: number): string {
  const m = Math.floor(seconds / 60);
  const s = seconds % 60;
  return `${m}:${s.toString().padStart(2, "0")}`;
}

function EventDetail() {
  const { id } = useParams<{ id: string }>();
  const [event, setEvent] = useState<Event | null>(null);
  const [error, setError] = useState("");
  const [now, setNow] = useState(() => Date.now());

  useEffect(() => {
    if (!id) return;
    api
      .get<Event>(`/events/${id}`)
      .then(setEvent)
      .catch((err) => setError(err.message));
  }, [id]);

  // Tick every second for countdown
  useEffect(() => {
    const timer = setInterval(() => setNow(Date.now()), 1000);
    return () => clearInterval(timer);
  }, []);

  if (error) return <p style={{ color: "red" }}>{error}</p>;
  if (!event) return <p>Loading...</p>;

  const ticketingStart = new Date(event.ticketing_starts_at).getTime();
  const ticketingEnd = new Date(event.ticketing_ends_at).getTime();
  const secondsUntilStart = Math.max(0, Math.ceil((ticketingStart - now) / 1000));
  const isOpen = now >= ticketingStart && now < ticketingEnd;
  const isClosed = now >= ticketingEnd;
  const showCountdown = secondsUntilStart > 0 && secondsUntilStart <= 600; // 10 minutes

  const handleBuyTickets = () => {
    window.open(
      `/events/${event.id}/tickets`,
      "ticketing",
      "width=500,height=600,scrollbars=yes"
    );
  };

  return (
    <div>
      <Link to="/">&larr; Back to events</Link>
      <h1>{event.name}</h1>
      <p>{event.venue}</p>
      <p>Event: {new Date(event.starts_at).toLocaleString("ja-JP")}</p>
      <p>
        Ticketing: {new Date(event.ticketing_starts_at).toLocaleString("ja-JP")}
        {" ~ "}
        {new Date(event.ticketing_ends_at).toLocaleString("ja-JP")}
      </p>

      {showCountdown && (
        <p style={{ fontSize: "2rem", fontWeight: "bold" }}>
          {formatCountdown(secondsUntilStart)}
        </p>
      )}

      {isClosed && <p>Ticketing has ended.</p>}
      {!isOpen && !isClosed && !showCountdown && <p>Ticketing has not started yet.</p>}

      <button onClick={handleBuyTickets} disabled={!isOpen}>
        {isOpen ? "Buy Tickets" : isClosed ? "Closed" : "Not Yet"}
      </button>
    </div>
  );
}

export default EventDetail;
