<script lang="ts">
  import { navigate } from "svelte-tiny-router";

  async function onclick() {
    const res = await fetch("/api/sessions/create", {
      method: "POST",
    });
    if (res.status !== 200) {
      alert("Failed to create new session. Please refresh and try again.");
    }
    const { uuid } = await res.json();
    navigate("/session/" + uuid);
  }
</script>

<div>
  <button {onclick}>Create New Voting Session</button>
</div>
