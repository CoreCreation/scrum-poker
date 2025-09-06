import { useLocation, useRoute } from "preact-iso";

import Panel from "./poker/panel";
import VoteList from "./poker/voteList";
import { useEffect, useRef } from "preact/hooks";
import { useState } from "react";
import EditNameModal from "./poker/editNameModal";

export default function Poker() {
  const {
    params: { id },
  } = useRoute();

  const { route } = useLocation();

  const [status, setStatus] = useState("Connecting..");
  const [voteData, setVoteData] = useState(null);
  const wsRef = useRef(null);
  const userName = useState(
    localStorage.getItem("scrum-poker-username") || null
  );

  // UI State
  const [editNameOpen, setEditNameOpen] = useState(false);

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
          body: "",
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
        setVoteData(msg.userData);
        console.log("Got Data", msg);
      } catch {
        console.log("Unable to parse JSON from WebSocket");
      }
    };

    return () => {
      wsRef.current = null;
      ws.close();
    };
  }, []);

  function editName(newName) {
    wsRef.current.send(
      JSON.stringify({
        type: "SetName",
        body: newName,
      })
    );
    setEditNameOpen(false);
  }

  function sendVote(number) {
    wsRef.current.send(
      JSON.stringify({
        type: "CastVote",
        body: String(number),
      })
    );
  }

  return (
    <div>
      Connection Status: {status} <br />
      Session ID: {id} <br />
      Username: {userName}
      <EditNameModal open={editNameOpen} save={editName} />
      <Panel options={[1, 2, 3, 5, 8, 12]} sendVote={sendVote} />
      <button onClick={() => setEditNameOpen((prev) => !prev)}>
        Edit Name
      </button>
      <button>Edit Vote Options</button>
      <button>Clear Votes</button>
      <VoteList data={voteData} />
    </div>
  );
}
