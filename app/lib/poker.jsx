import { useLocation, useRoute } from "preact-iso";

import Panel from "./poker/panel";
import VoteList from "./poker/voteList";
import { useEffect, useRef, useState } from "preact/hooks";
import EditNameModal from "./poker/editNameModal";
import EditVoteOptionsModal from "./poker/editVoteOptionsModal";
import JSConfetti from "js-confetti";

import ThemeToggle from "./components/themeToggle";

export default function Poker() {
  const {
    params: { id },
  } = useRoute();

  const { route } = useLocation();

  const userId = useRef(null);
  const [status, setStatus] = useState("Connecting..");
  const [votesVisible, setVotesVisible] = useState(false);
  const [voteOptions, setVoteOptions] = useState([]);
  const [voteData, setVoteData] = useState(null);
  const wsRef = useRef(null);
  const wsStopped = useRef(false);
  const confettiRef = useRef(null);
  const dropdownRef = useRef(null);
  const [username, setUserName] = useState(
    localStorage.getItem("scrum-poker-username") || null
  );
  const [isVoting, setIsVoting] = useState(() => {
    const curr = localStorage.getItem("scrum-poker-is-voting");
    return curr === null ? true : curr === "true";
  });
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

    startWebSocket();

    return () => {
      wsStopped.current = true;
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, []);

  function startWebSocket() {
    if (wsStopped.current === true) {
      return;
    }

    // Get WebSocket Connection
    const url = new URL("/api/sessions/" + id + "/join", window.location.href);
    url.protocol = url.protocol.replace("http", "ws");
    console.log("Connecting to", url.href);

    const ws = new WebSocket(url);
    wsRef.current = ws;

    ws.onopen = () => {
      setStatus("Open");
      const initBody = {
        type: "Init",
      };
      if (username) initBody.username = username;
      if (userId.current) initBody.body = userId.current;
      if (lastVote) initBody.vote = lastVote;
      if (isVoting === false) initBody.active = false;
      ws.send(JSON.stringify(initBody));
    };
    ws.onclose = () => {
      ws.close();
      wsRef.current = null;
      startWebSocket();
      setStatus("Connecting...");
    };
    ws.onerror = () => {
      wsStopped.current = true;
      alert("Error with Connection. Please Refresh to Reconnect");
      setStatus("Error");
    };
    ws.onmessage = (e) => {
      try {
        const msg = JSON.parse(e.data);
        console.log(msg);
        userId.current = msg.userId;
        setVotesVisible(msg.votesVisible);
        setVoteData(msg.userData);
        if (msg.username.length && username !== msg.username) {
          persistUsername(msg.username);
        }
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
  }

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
        type: "UpdateData",
        username: newName,
      })
    );
    persistUsername(newName);
    setEditNameOpen(false);
  }

  function persistUsername(newName) {
    setUserName(newName);
    localStorage.setItem("scrum-poker-username", newName);
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
        type: "UpdateData",
        vote: number,
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

  function leaveVote() {
    wsRef.current.send(
      JSON.stringify({
        type: "UpdateData",
        active: false,
      })
    );
    setLastVote(null);
    setIsVoting(false);
    dropdownRef.current.open = false;
    localStorage.setItem("scrum-poker-is-voting", false);
  }

  function joinVote() {
    wsRef.current.send(
      JSON.stringify({
        type: "UpdateData",
        active: true,
      })
    );
    setIsVoting(true);
    dropdownRef.current.open = false;
    localStorage.setItem("scrum-poker-is-voting", true);
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
              <details class="dropdown" ref={dropdownRef}>
                <summary>Actions</summary>
                <ul dir="rtl">
                  <li>
                    <a
                      onClick={() => {
                        setEditNameOpen((prev) => !prev);
                        dropdownRef.current.open = false;
                      }}
                    >
                      Edit Username
                    </a>
                  </li>
                  <li>
                    <a
                      onClick={() => {
                        setEditVotesOpen(true);
                        dropdownRef.current.open = false;
                      }}
                    >
                      Edit Vote Options
                    </a>
                  </li>
                  <li>
                    {isVoting ? (
                      <a onClick={leaveVote}>Leave Vote</a>
                    ) : (
                      <a onClick={joinVote}>Join Vote</a>
                    )}
                  </li>
                </ul>
              </details>
            </li>
            <ThemeToggle />
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
            name={username}
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
            disabled={!isVoting}
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
