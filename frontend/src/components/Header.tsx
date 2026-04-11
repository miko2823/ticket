import { signOut } from "firebase/auth";
import { auth } from "../auth/firebase";
import { useAuth } from "../auth/AuthContext";

function Header() {
  const { user } = useAuth();

  if (!user) return null;

  return (
    <header>
      <span>{user.email}</span>
      <button onClick={() => signOut(auth)}>Sign Out</button>
    </header>
  );
}

export default Header;
