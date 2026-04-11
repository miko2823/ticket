import { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";
import { api } from "../api/client";

interface Event {
  id: string;
  name: string;
  venue: string;
  starts_at: string;
}

function EventDetail() {
  const { id } = useParams<{ id: string }>();
  const [event, setEvent] = useState<Event | null>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!id) return;
    api
      .get<Event>(`/events/${id}`)
      .then(setEvent)
      .catch((err) => setError(err.message));
  }, [id]);

  if (error) return <p style={{ color: "red" }}>{error}</p>;
  if (!event) return <p>Loading...</p>;

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
      <p>{new Date(event.starts_at).toLocaleString("ja-JP")}</p>
      <button onClick={handleBuyTickets}>Buy Tickets</button>
    </div>
  );
}

export default EventDetail;
