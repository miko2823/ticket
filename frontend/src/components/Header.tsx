import { signOut } from "firebase/auth";
import { useNavigate, Link } from "react-router-dom";
import { auth } from "../auth/firebase";
import { useAuth } from "../auth/AuthContext";

function Header() {
  const { user } = useAuth();
  const navigate = useNavigate();

  if (!user) return null;

  const handleSignOut = async () => {
    await signOut(auth);
    navigate("/");
  };

  return (
    <header>
      <span>{user.email}</span>
      <Link to="/my-tickets">My Bookings</Link>
      <button onClick={handleSignOut}>Sign Out</button>
    </header>
  );
}

export default Header;
