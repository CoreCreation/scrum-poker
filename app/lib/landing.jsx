import { useLocation } from "preact-iso";

export default function Landing() {
  const { route } = useLocation();

  async function onclick() {
    const res = await fetch("/api/sessions/create", {
      method: "POST",
    });
    if (res.status !== 200) {
      alert("Failed to create new session. Please refresh and try again.");
    }
    const { uuid } = await res.json();
    route("/session/" + uuid);
  }

  return (
    <div class="landing-page">
      <button onClick={onclick}>Create New Voting Session</button>
    </div>
  );
}
