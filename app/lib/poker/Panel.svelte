<script lang="ts">
  import VoteButton from "./panel/VoteButton.svelte";

  let { options } = $props();
  let str = $derived(options.join(", "));
</script>

<div>
  {#each options as option}
    <VoteButton {option} />
  {/each}
</div>
