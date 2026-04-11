import { useEffect, useState } from "react";
import { api } from "../api/client";

function Home() {
  const [uid, setUid] = useState<string | null>(null);

  useEffect(() => {
    api.get<{ uid: string }>("/me").then((data) => setUid(data.uid)).catch(console.error);
  }, []);

  return (
    <div>
      <h1>SturdyTicket</h1>
      {uid && <p>Authenticated as: {uid}</p>}
      <p>Events will be listed here.</p>
    </div>
  );
}

export default Home;
