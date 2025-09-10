import { useLocation, useRoute } from "preact-iso";

import Panel from "./poker/panel";
import VoteList from "./poker/voteList";
import { useEffect, useRef, useState } from "preact/hooks";
import EditNameModal from "./poker/editNameModal";
import EditVoteOptionsModal from "./poker/editVoteOptionsModal";
import JSConfetti from "js-confetti";
import { v4 as uuid } from "uuid";

import ThemeToggle from "./components/themeToggle";

export default function Poker() {
  const {
    params: { id },
  } = useRoute();

  const { route } = useLocation();

  const clientId = useRef(null);
  const [status, setStatus] = useState("Connecting..");
  const [votesVisible, setVotesVisible] = useState(false);
  const [voteOptions, setVoteOptions] = useState([]);
  const [clientData, setClientData] = useState(null);
  const wsStateRef = useRef({});
  const confettiRef = useRef(null);
  const dropdownRef = useRef(null);
  const [username, setUserName] = useState(
    localStorage.getItem("scrum-poker-username") || null
  );
  const [active, setActive] = useState(() => {
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
      const votes = clientData.map((c) => c.vote).filter((c) => c !== -1);
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

    const key = "scrum-poker-client-id";
    let cid = localStorage.getItem(key);
    if (!cid) {
      cid = uuid();
      localStorage.setItem(key, cid);
    }
    clientId.current = cid;

    startWebSocket();

    return () => {
      wsStateRef.current.stopped = true;
      if (wsStateRef.current.ws) {
        wsStateRef.current.ws.close();
      }
    };
  }, []);

  function startWebSocket() {
    if (wsStateRef.current.stopped === true) {
      return;
    }

    // Get WebSocket Connection
    const url = new URL(
      "/api/sessions/" + id + "/join/" + clientId.current,
      window.location.href
    );
    url.protocol = url.protocol.replace("http", "ws");
    console.log("Connecting to", url.href, "with client ID", clientId.current);

    const ws = new WebSocket(url);
    wsStateRef.current.ws = ws;

    ws.onopen = () => {
      setStatus("Open");
    };
    ws.onclose = () => {
      ws.close();
      wsStateRef.current.ws = null;
      startWebSocket();
      setStatus("Connecting...");
    };
    ws.onerror = () => {
      wsStateRef.current.stopped = true;
      alert("Error with Connection. Please Refresh to Reconnect");
      setStatus("Error");
    };
    ws.onmessage = (e) => {
      try {
        const msg = JSON.parse(e.data);
        console.log(msg);
        if (msg.type === "Init") {
          const obj = {
            type: "UpdateData",
          };
          if (username) obj.username = username;
          if (lastVote) obj.vote = lastVote;
          if (active === false) obj.active = active;
          ws.send(JSON.stringify(obj));
        } else {
          persistUsername(msg.username);
          setLastVote(msg.vote);
          setVotesVisible(msg.votesVisible);
          setActive(msg.active);
          setClientData(msg.clientData);
          const numbers = msg.voteOptions
            .split(",")
            .map((s) => parseInt(s.trim()));
          if (
            numbers.reduce((p, c) => (p &&= c > 0 && !Number.isNaN(c)), true) &&
            !isEqual(voteOptions, numbers)
          ) {
            setVoteOptions(numbers);
          }
        }
      } catch (e) {
        console.log("Unable to parse JSON from WebSocket", e);
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
    wsStateRef.current.ws.send(
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
    wsStateRef.current.ws.send(
      JSON.stringify({
        type: "SetOptions",
        body: str,
      })
    );
    setEditVotesOpen(false);
  }

  function sendVote(number) {
    wsStateRef.current.ws.send(
      JSON.stringify({
        type: "UpdateData",
        vote: number,
      })
    );
    setLastVote(number);
  }

  function clearVotes() {
    wsStateRef.current.ws.send(
      JSON.stringify({
        type: "ClearVotes",
        body: "",
      })
    );
  }

  function showVotes() {
    wsStateRef.current.ws.send(
      JSON.stringify({
        type: "ShowVotes",
        body: "",
      })
    );
  }

  function leaveVote() {
    wsStateRef.current.ws.send(
      JSON.stringify({
        type: "UpdateData",
        active: false,
      })
    );
    setLastVote(null);
    setActive(false);
    dropdownRef.current.open = false;
    localStorage.setItem("scrum-poker-is-voting", false);
  }

  function joinVote() {
    wsStateRef.current.ws.send(
      JSON.stringify({
        type: "UpdateData",
        active: true,
      })
    );
    setActive(true);
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
                    {active ? (
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
      {!clientData ? (
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
            disabled={!active}
          />
          <VoteList data={clientData} votesVisible={votesVisible} />
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
