import { useEffect, useState } from "react";
import { api } from "../api/client";

interface Booking {
  id: string;
  ticket_id: string;
  status: string;
  created_at: string;
}

function MyTickets() {
  const [bookings, setBookings] = useState<Booking[]>([]);

  const load = () => {
    api.get<Booking[]>("/bookings/me").then(setBookings).catch(() => {});
  };

  useEffect(load, []);

  const handleCancel = async (id: string) => {
    if (!confirm("Cancel this booking?")) return;
    try {
      await api.del(`/bookings/${id}`);
      load();
    } catch (err) {
      alert(err instanceof Error ? err.message : "Failed to cancel");
    }
  };

  return (
    <div>
      <h1>My Tickets</h1>
      {bookings.length === 0 && <p>No bookings yet.</p>}
      <ul>
        {bookings.map((b) => (
          <li key={b.id}>
            Booking: {b.id.slice(0, 8)}... — <strong>{b.status}</strong> —{" "}
            {new Date(b.created_at).toLocaleString("ja-JP")}
            {(b.status === "confirmed" || b.status === "pending") && (
              <button onClick={() => handleCancel(b.id)} style={{ marginLeft: 8 }}>
                Cancel
              </button>
            )}
          </li>
        ))}
      </ul>
    </div>
  );
}

export default MyTickets;
