import { useLocation, useRoute } from "preact-iso";

import Panel from "./poker/panel";
import VoteList from "./poker/voteList";
import { useEffect, useRef, useState } from "preact/hooks";
import EditNameModal from "./poker/editNameModal";
import EditVoteOptionsModal from "./poker/editVoteOptionsModal";
import JSConfetti from "js-confetti";

export default function Poker() {
  const {
    params: { id },
  } = useRoute();

  const { route } = useLocation();

  const [status, setStatus] = useState("Connecting..");
  const [votesVisible, setVotesVisible] = useState(false);
  const [voteOptions, setVoteOptions] = useState([]);
  const [voteData, setVoteData] = useState(null);
  const wsRef = useRef(null);
  const confettiRef = useRef(null);
  const [userName, setUserName] = useState(
    localStorage.getItem("scrum-poker-username") || null
  );
  const [lastVote, setLastVote] = useState(null);

  // UI State
  const [editNameOpen, setEditNameOpen] = useState(false);
  const [editVotesOpen, setEditVotesOpen] = useState(false);

  useEffect(() => {
    confettiRef.current = new JSConfetti();

    return () => confettiRef.current?.clearCanvas();
  }, []);

  useEffect(() => {
    if (votesVisible) {
      const votes = voteData.map((c) => c.vote).filter((c) => c !== -1);
      const first = votes[0];
      if (
        votes.length > 1 &&
        votes.reduce((p, c) => (p &&= c === first), true)
      ) {
        confettiRef.current.addConfetti();
      }
    }
    setLastVote(null);
  }, [votesVisible]);

  useEffect(async () => {
    let res = await fetch("/api/sessions/" + id);
    if (res.status !== 200) {
      alert("Session no longer valid, please create a new one.");
      return route("/");
    }

    // Get WebSocket Connection
    const url = new URL("/api/sessions/" + id + "/join", window.location.href);
    url.protocol = url.protocol.replace("http", "ws");
    console.log("Connecting to", url.href);

    const ws = new WebSocket(url);
    wsRef.current = ws;

    ws.onopen = () => {
      setStatus("Open");
      ws.send(
        JSON.stringify({
          type: "Init",
          body: userName ? userName : "",
        })
      );
    };
    ws.onclose = () => {
      alert("Connection Closed. Please Refresh to Reconnect");
      setStatus("Closed");
    };
    ws.onerror = () => {
      alert("Error with Connection. Please Refresh to Reconnect");
      setStatus("Error");
    };
    ws.onmessage = (e) => {
      try {
        const msg = JSON.parse(e.data);
        setVotesVisible(msg.votesVisible);
        setVoteData(msg.userData);
        const numbers = msg.voteOptions
          .split(",")
          .map((s) => parseInt(s.trim()));
        if (
          numbers.reduce((p, c) => (p &&= c > 0 && !Number.isNaN(c)), true) &&
          !isEqual(voteOptions, numbers)
        ) {
          setVoteOptions(numbers);
        }
      } catch {
        console.log("Unable to parse JSON from WebSocket");
      }
    };

    return () => {
      wsRef.current = null;
      ws.close();
    };
  }, []);

  function isEqual(a1, a2) {
    if (a1.length !== a2.length) return false;
    for (let i = 0; i < a1.length; i++) {
      if (a1[i] !== a2[i]) return false;
    }
    return true;
  }

  function editName(newName) {
    wsRef.current.send(
      JSON.stringify({
        type: "SetName",
        body: newName,
      })
    );
    setUserName(newName);
    localStorage.setItem("scrum-poker-username", newName);
    setEditNameOpen(false);
  }

  function editVotes(str) {
    wsRef.current.send(
      JSON.stringify({
        type: "SetOptions",
        body: str,
      })
    );
    setEditVotesOpen(false);
  }

  function sendVote(number) {
    wsRef.current.send(
      JSON.stringify({
        type: "CastVote",
        body: String(number),
      })
    );
    setLastVote(number);
  }

  function clearVotes() {
    wsRef.current.send(
      JSON.stringify({
        type: "ClearVotes",
        body: "",
      })
    );
  }

  function showVotes() {
    wsRef.current.send(
      JSON.stringify({
        type: "ShowVotes",
        body: "",
      })
    );
  }

  return (
    <div class="container poker">
      <header>
        <nav>
          <ul>
            <li>
              <strong>&#123;TeamNameGoesHere&#125; Poker</strong>
            </li>
          </ul>
          <ul>
            <li>
              <details class="dropdown">
                <summary>Actions</summary>
                <ul dir="rtl">
                  <li>
                    <a onClick={() => setEditNameOpen((prev) => !prev)}>
                      Edit Username
                    </a>
                  </li>
                  <li>
                    <a onClick={() => setEditVotesOpen(true)}>
                      Edit Vote Options
                    </a>
                  </li>
                </ul>
              </details>
            </li>
          </ul>
        </nav>
      </header>
      {!voteData ? (
        <main>
          <div aria-busy="true"></div>
        </main>
      ) : (
        <main>
          <EditNameModal
            open={editNameOpen}
            setOpen={setEditNameOpen}
            save={editName}
            name={userName}
          />
          <EditVoteOptionsModal
            open={editVotesOpen}
            setOpen={setEditVotesOpen}
            save={editVotes}
            current={voteOptions}
          />
          <Panel
            options={voteOptions}
            sendVote={sendVote}
            lastVote={lastVote}
          />
          <VoteList data={voteData} votesVisible={votesVisible} />
          {votesVisible ? (
            <button class="secondary" onClick={clearVotes}>
              Clear Votes
            </button>
          ) : (
            <button onClick={showVotes}>Show Votes</button>
          )}
        </main>
      )}
      <footer>
        <span>Connection Status: {status}</span>
        <span>
          Session ID: {id} <br />
        </span>
      </footer>
    </div>
  );
}
