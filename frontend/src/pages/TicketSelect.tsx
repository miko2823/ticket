import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { api } from "../api/client";
import styles from "./TicketSelect.module.css";

interface Ticket {
  id: string;
  event_id: string;
  seat_label: string;
  price_jpy: number;
  status: string;
}

function TicketSelect() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [selected, setSelected] = useState<Ticket | null>(null);
  const [reserving, setReserving] = useState(false);

  useEffect(() => {
    if (!id) return;
    api
      .get<Ticket[]>(`/events/${id}/tickets`)
      .then(setTickets)
      .catch((err) => alert(err.message));
  }, [id]);

  const handleSelect = (ticket: Ticket) => {
    if (ticket.status !== "available") return;
    setSelected(ticket);
  };

  const handleProceed = async () => {
    if (!selected) return;
    setReserving(true);
    try {
      await api.post(`/tickets/${selected.id}/reserve`, {});
      navigate(`/checkout?ticketId=${selected.id}&eventId=${id}`);
    } catch (err) {
      const msg = err instanceof Error ? err.message : "Failed to reserve seat";
      alert(`This seat is no longer available: ${msg}\nPlease select another seat.`);
      setSelected(null);
      // Refresh tickets to show updated availability
      if (id) {
        api.get<Ticket[]>(`/events/${id}/tickets`).then(setTickets).catch(() => {});
      }
    } finally {
      setReserving(false);
    }
  };

  return (
    <div className={styles.container}>
      <h2>Select a Seat</h2>
      <div className={styles.legend}>
        <span className={`${styles.seat} ${styles.available}`} /> Available
        <span className={`${styles.seat} ${styles.reserved}`} /> Reserved
        <span className={`${styles.seat} ${styles.selectedSeat}`} /> Selected
      </div>
      <div className={styles.grid}>
        {tickets.map((t) => (
          <button
            key={t.id}
            className={`${styles.seat} ${
              selected?.id === t.id
                ? styles.selectedSeat
                : t.status === "available"
                  ? styles.available
                  : styles.reserved
            }`}
            onClick={() => handleSelect(t)}
            disabled={t.status !== "available"}
            title={`${t.seat_label} — ¥${t.price_jpy.toLocaleString()}`}
          >
            {t.seat_label}
          </button>
        ))}
      </div>
      {selected && (
        <div className={styles.selectedInfo}>
          <p>
            Seat: <strong>{selected.seat_label}</strong> — ¥
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
