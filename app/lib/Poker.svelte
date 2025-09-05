<script lang="ts">
  import { onMount } from "svelte";
  import Panel from "./poker/Panel.svelte";
  import VoteList from "./poker/VoteList.svelte";
  import { navigate } from "svelte-tiny-router";

  let { id } = $props();

  onMount(async () => {
    console.log(id);
    let res = await fetch("/api/sessions/" + id);
    if (res.status !== 200) {
      alert("Session no longer valid, please create a new one.");
      navigate("/");
    }
  });
</script>

<div>
  {id}
  <Panel options={[1, 2, 3, 5, 8, 12]} />
  <button>Edit Vote Options</button>
  <button>Clear Votes</button>
  <VoteList />
</div>
