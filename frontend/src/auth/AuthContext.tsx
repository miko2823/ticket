import { createContext, useContext, useState, useEffect, type ReactNode } from "react";

interface AuthState {
  uid: string | null;
  loading: boolean;
}

const AuthContext = createContext<AuthState>({ uid: null, loading: true });

export function AuthProvider({ children }: { children: ReactNode }) {
  const [state, setState] = useState<AuthState>({ uid: null, loading: true });

  useEffect(() => {
    // TODO: initialize Firebase Auth and listen for auth state changes
    setState({ uid: null, loading: false });
  }, []);

  return <AuthContext.Provider value={state}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  return useContext(AuthContext);
}
