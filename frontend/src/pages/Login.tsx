import { useState } from "react";
import {
  signInWithEmailAndPassword,
  createUserWithEmailAndPassword,
} from "firebase/auth";
import { auth } from "../auth/firebase";

function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [isSignUp, setIsSignUp] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    try {
      if (isSignUp) {
        await createUserWithEmailAndPassword(auth, email, password);
      } else {
        await signInWithEmailAndPassword(auth, email, password);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Authentication failed");
    }
  };

  return (
    <div>
      <h1>SturdyTicket</h1>
      <h2>{isSignUp ? "Sign Up" : "Sign In"}</h2>
      <form onSubmit={handleSubmit}>
        <div>
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>
        <div>
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        {error && <p style={{ color: "red" }}>{error}</p>}
        <button type="submit">{isSignUp ? "Sign Up" : "Sign In"}</button>
      </form>
      <p>
        {isSignUp ? "Already have an account?" : "Don't have an account?"}{" "}
        <button type="button" onClick={() => setIsSignUp(!isSignUp)}>
          {isSignUp ? "Sign In" : "Sign Up"}
        </button>
      </p>
    </div>
  );
}

export default Login;
