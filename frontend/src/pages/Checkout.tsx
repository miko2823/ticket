import { useEffect, useState } from "react";
import { useSearchParams } from "react-router-dom";
import { api } from "../api/client";

interface Ticket {
  id: string;
  seat_label: string;
  price_jpy: number;
  status: string;
  reserved_until?: string;
}

interface BookingResponse {
  id: string;
  status: string;
}

function Checkout() {
  const [searchParams] = useSearchParams();
  const ticketId = searchParams.get("ticketId");
  const [ticket, setTicket] = useState<Ticket | null>(null);
  const [booking, setBooking] = useState<BookingResponse | null>(null);
  const [paying, setPaying] = useState(false);
  const [error, setError] = useState("");
  const [timeLeft, setTimeLeft] = useState<number | null>(null);

  useEffect(() => {
    if (!ticketId) return;
    api
      .get<Ticket>(`/tickets/${ticketId}`)
      .catch(() => null)
      .then((data) => {
        // Fallback: fetch from event tickets if direct endpoint doesn't exist
        if (!data) return;
        setTicket(data);
      });
  }, [ticketId]);

  // Countdown timer for reservation expiry
  useEffect(() => {
    if (!ticket?.reserved_until) return;
    const expiresAt = new Date(ticket.reserved_until).getTime();

    const interval = setInterval(() => {
      const remaining = Math.max(0, Math.floor((expiresAt - Date.now()) / 1000));
      setTimeLeft(remaining);
      if (remaining === 0) {
        setError("Reservation expired. Please go back and select a seat again.");
        clearInterval(interval);
      }
    }, 1000);

    return () => clearInterval(interval);
  }, [ticket?.reserved_until]);

  const handlePay = async () => {
    if (!ticketId) return;
    setPaying(true);
    setError("");
    try {
      const result = await api.post<BookingResponse>("/bookings", {
        ticket_id: ticketId,
      });
      setBooking(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Payment failed");
    } finally {
      setPaying(false);
    }
  };

  if (!ticketId) return <p>No ticket selected.</p>;

  if (booking) {
    return (
      <div>
        <h1>Booking Confirmed</h1>
        <p>Booking ID: {booking.id}</p>
        <p>Status: {booking.status}</p>
        <p>You can close this window.</p>
      </div>
    );
  }

  const formatTime = (seconds: number) => {
    const m = Math.floor(seconds / 60);
    const s = seconds % 60;
    return `${m}:${s.toString().padStart(2, "0")}`;
  };

  return (
    <div>
      <h1>Checkout</h1>
      {ticket && (
        <div>
          <p>
            Seat: <strong>{ticket.seat_label}</strong>
          </p>
          <p>Price: ¥{ticket.price_jpy.toLocaleString()}</p>
        </div>
      )}
      {timeLeft !== null && timeLeft > 0 && (
        <p>Reservation expires in: {formatTime(timeLeft)}</p>
      )}
      {error && <p style={{ color: "red" }}>{error}</p>}
      {!error && (
        <button onClick={handlePay} disabled={paying}>
          {paying ? "Processing..." : "Pay"}
        </button>
      )}
    </div>
  );
}

export default Checkout;
