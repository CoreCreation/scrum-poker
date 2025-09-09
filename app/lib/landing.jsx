import { useLocation } from "preact-iso";
import ThemeToggle from "./components/themeToggle";

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
    <div class="container landing">
      <header>
        <nav>
          <ul>
            <li>
              <strong>&#123;TeamNameGoesHere&#125; Poker</strong>
            </li>
          </ul>
          <ul>
            <ThemeToggle />
          </ul>
        </nav>
      </header>
      <main>
        <button onClick={onclick}>Create New Voting Session</button>
      </main>
    </div>
  );
}
