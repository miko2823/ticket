import { Routes, Route } from "react-router-dom";
import { useAuth } from "./auth/AuthContext";
import Header from "./components/Header";
import Home from "./pages/Home";
import EventDetail from "./pages/EventDetail";
import TicketSelect from "./pages/TicketSelect";
import Checkout from "./pages/Checkout";
import MyTickets from "./pages/MyTickets";
import Login from "./pages/Login";

function App() {
  const { user, loading } = useAuth();

  if (loading) return <p>Loading...</p>;
  if (!user) return <Login />;

  return (
    <>
      <Header />
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/events/:id" element={<EventDetail />} />
        <Route path="/events/:id/tickets" element={<TicketSelect />} />
        <Route path="/checkout" element={<Checkout />} />
        <Route path="/my-tickets" element={<MyTickets />} />
      </Routes>
    </>
  );
}

export default App;
