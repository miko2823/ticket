import { useEffect, useRef, useState, useCallback } from "react";
import { api } from "../api/client";
import styles from "./SeatMap.module.css";

interface Section {
  id: string;
  label: string;
  color: string;
}

interface SeatData {
  ticket_id: string;
  label: string;
  section: string;
  x: number;
  y: number;
  r: number;
  price_jpy: number;
  status: string;
}

interface SeatMapData {
  event_id: string;
  layout: {
    canvas: { width: number; height: number };
    stage: { x: number; y: number; width: number; height: number; label: string };
    sections: Section[];
  };
  seats: SeatData[];
}

interface SeatMapProps {
  eventId: string;
  selectedTicketId: string | null;
  onSeatSelect: (seat: SeatData) => void;
}

const RESERVED_COLOR = "#BDBDBD";
const SELECTED_COLOR = "#1565C0";

function SeatMap({ eventId, selectedTicketId, onSeatSelect }: SeatMapProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [data, setData] = useState<SeatMapData | null>(null);
  const [tooltip, setTooltip] = useState<{ x: number; y: number; text: string } | null>(null);

  useEffect(() => {
    api
      .get<SeatMapData>(`/events/${eventId}/seatmap`)
      .then(setData)
      .catch((err) => alert(err.message));
  }, [eventId]);

  const getSectionColor = useCallback(
    (sectionId: string): string => {
      if (!data) return "#999";
      const sec = data.layout.sections.find((s) => s.id === sectionId);
      return sec?.color ?? "#999";
    },
    [data]
  );

  const draw = useCallback(() => {
    const canvas = canvasRef.current;
    if (!canvas || !data) return;
    const ctx = canvas.getContext("2d");
    if (!ctx) return;

    const { layout, seats } = data;
    const dpr = window.devicePixelRatio || 1;

    canvas.width = layout.canvas.width * dpr;
    canvas.height = layout.canvas.height * dpr;
    canvas.style.aspectRatio = `${layout.canvas.width} / ${layout.canvas.height}`;
    ctx.scale(dpr, dpr);

    // Background
    ctx.fillStyle = "#F5F5F5";
    ctx.fillRect(0, 0, layout.canvas.width, layout.canvas.height);

    // Stage
    const st = layout.stage;
    ctx.fillStyle = "#333";
    ctx.fillRect(st.x, st.y, st.width, st.height);
    ctx.fillStyle = "#FFF";
    ctx.font = "bold 14px sans-serif";
    ctx.textAlign = "center";
    ctx.textBaseline = "middle";
    ctx.fillText(st.label, st.x + st.width / 2, st.y + st.height / 2);

    // Seats
    for (const seat of seats) {
      const isSelected = seat.ticket_id === selectedTicketId;
      const isAvailable = seat.status === "available";

      // Fill color
      if (isSelected) {
        ctx.fillStyle = SELECTED_COLOR;
      } else if (isAvailable) {
        ctx.fillStyle = getSectionColor(seat.section);
      } else {
        ctx.fillStyle = RESERVED_COLOR;
      }

      ctx.beginPath();
      ctx.arc(seat.x, seat.y, seat.r, 0, Math.PI * 2);
      ctx.fill();

      // Border
      ctx.strokeStyle = isSelected ? "#0D47A1" : "rgba(0,0,0,0.2)";
      ctx.lineWidth = isSelected ? 3 : 1;
      ctx.stroke();

      // Label
      ctx.fillStyle = isSelected || !isAvailable ? "#FFF" : "#000";
      ctx.font = `${Math.max(9, seat.r - 4)}px sans-serif`;
      ctx.textAlign = "center";
      ctx.textBaseline = "middle";
      ctx.fillText(seat.label, seat.x, seat.y);
    }
  }, [data, selectedTicketId, getSectionColor]);

  useEffect(() => {
    draw();
  }, [draw]);

  const toCanvasCoords = (e: React.MouseEvent<HTMLCanvasElement>) => {
    const canvas = canvasRef.current;
    if (!canvas || !data) return null;
    const rect = canvas.getBoundingClientRect();
    const scaleX = data.layout.canvas.width / rect.width;
    const scaleY = data.layout.canvas.height / rect.height;
    return {
      x: (e.clientX - rect.left) * scaleX,
      y: (e.clientY - rect.top) * scaleY,
    };
  };

  const findSeatAt = (cx: number, cy: number): SeatData | null => {
    if (!data) return null;
    for (const seat of data.seats) {
      const dx = cx - seat.x;
      const dy = cy - seat.y;
      if (dx * dx + dy * dy <= seat.r * seat.r) return seat;
    }
    return null;
  };

  const handleClick = (e: React.MouseEvent<HTMLCanvasElement>) => {
    const pt = toCanvasCoords(e);
    if (!pt) return;
    const seat = findSeatAt(pt.x, pt.y);
    if (seat && seat.status === "available") {
      onSeatSelect(seat);
    }
  };

  const handleMouseMove = (e: React.MouseEvent<HTMLCanvasElement>) => {
    const pt = toCanvasCoords(e);
    if (!pt) {
      setTooltip(null);
      return;
    }
    const seat = findSeatAt(pt.x, pt.y);
    if (seat) {
      const canvas = canvasRef.current!;
      const rect = canvas.getBoundingClientRect();
      setTooltip({
        x: e.clientX - rect.left,
        y: e.clientY - rect.top,
        text: `${seat.label} — ¥${seat.price_jpy.toLocaleString()} (${seat.status})`,
      });
    } else {
      setTooltip(null);
    }
  };

  if (!data) return <p>Loading seat map...</p>;

  return (
    <div className={styles.container}>
      <canvas
        ref={canvasRef}
        className={styles.canvas}
        onClick={handleClick}
        onMouseMove={handleMouseMove}
        onMouseLeave={() => setTooltip(null)}
      />
      {tooltip && (
        <div className={styles.tooltip} style={{ left: tooltip.x, top: tooltip.y }}>
          {tooltip.text}
        </div>
      )}
    </div>
  );
}

export default SeatMap;
export type { SeatData, Section };
