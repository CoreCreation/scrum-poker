<script lang="ts">
  import Landing from "./lib/Landing.svelte";
  import Poker from "./lib/Poker.svelte";

  import { Router, Route } from "svelte-tiny-router";

  export let url = "";
</script>

<Router {url}>
  <main>
    <Route path="/" component={Landing} />
    <Route path="/session/:id" component={Poker} />
  </main>
</Router>

<style>
</style>
