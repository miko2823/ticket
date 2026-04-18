import { useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { api } from "../api/client";
import { useSession } from "../hooks/useSession";
import SeatMap from "../components/SeatMap";
import type { SeatData } from "../components/SeatMap";
import styles from "./TicketSelect.module.css";

function TicketSelect() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [selected, setSelected] = useState<SeatData | null>(null);
  const [reserving, setReserving] = useState(false);
  const [key, setKey] = useState(0);

  const { sessionId, error: sessionError, isConnecting } = useSession(id ?? "");

  const handleSeatSelect = (seat: SeatData) => {
    setSelected(seat);
  };

  const handleProceed = async () => {
    if (!selected || !sessionId.current) return;
    setReserving(true);
    try {
      await api.post(`/tickets/${selected.ticket_id}/reserve`, {}, {
        "X-Session-ID": sessionId.current,
      });
      navigate(`/checkout?ticketId=${selected.ticket_id}&eventId=${id}`);
    } catch (err) {
      const msg = err instanceof Error ? err.message : "Failed to reserve seat";
      alert(`This seat is no longer available: ${msg}\nPlease select another seat.`);
      setSelected(null);
      setKey((k) => k + 1);
    } finally {
      setReserving(false);
    }
  };

  if (!id) return <p>No event selected.</p>;
  if (isConnecting) return <p>Connecting...</p>;
  if (sessionError) return <p style={{ color: "red" }}>{sessionError}</p>;

  return (
    <div className={styles.container}>
      <h2>Select a Seat</h2>
      <SeatMap
        key={key}
        eventId={id}
        selectedTicketId={selected?.ticket_id ?? null}
        onSeatSelect={handleSeatSelect}
      />
      {selected && (
        <div className={styles.selectedInfo}>
          <p>
            Seat: <strong>{selected.label}</strong> — ¥
            {selected.price_jpy.toLocaleString()}
          </p>
          <button onClick={handleProceed} disabled={reserving}>
            {reserving ? "Reserving..." : "Proceed to Checkout"}
          </button>
        </div>
      )}
    </div>
  );
}

export default TicketSelect;
