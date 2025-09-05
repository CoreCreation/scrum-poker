import { useLocation, useRoute } from "preact-iso";

import Panel from "./poker/panel";
import VoteList from "./poker/voteList";
import { useEffect } from "preact/hooks";

export default function Poker() {
  const {
    params: { id },
  } = useRoute();

  const { route } = useLocation();

  useEffect(async () => {
    let res = await fetch("/api/sessions/" + id);
    console.log(res);
    if (res.status !== 200) {
      alert("Session no longer valid, please create a new one.");
      route("/");
    }
  }, []);

  return (
    <div>
      Session ID: {id}
      <Panel options={[1, 2, 3, 5, 8, 12]} />
      <button>Edit Vote Options</button>
      <button>Clear Votes</button>
      <VoteList />
    </div>
  );
}
