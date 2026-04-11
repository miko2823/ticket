import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { api } from "../api/client";

interface Event {
  id: string;
  name: string;
  venue: string;
  starts_at: string;
  ticketing_starts_at: string;
  ticketing_ends_at: string;
}

function Home() {
  const [events, setEvents] = useState<Event[]>([]);
  const [error, setError] = useState("");

  useEffect(() => {
    api
      .get<Event[]>("/events")
      .then(setEvents)
      .catch((err) => setError(err.message));
  }, []);

  return (
    <div>
      <h1>SturdyTicket</h1>
      {error && <p style={{ color: "red" }}>{error}</p>}
      <ul>
        {events.map((e) => (
          <li key={e.id}>
            <Link to={`/events/${e.id}`}>
              <strong>{e.name}</strong>
            </Link>
            <br />
            {e.venue} — {new Date(e.starts_at).toLocaleDateString("ja-JP")}
          </li>
        ))}
      </ul>
    </div>
  );
}

export default Home;
