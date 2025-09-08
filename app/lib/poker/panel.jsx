import VoteButton from "./panel/voteButton";

export default function Panel({ options, sendVote, lastVote }) {
  return (
    <div class="vote-panel-options">
      {options.map((option) => (
        <VoteButton option={option} sendVote={sendVote} lastVote={lastVote} />
      ))}
    </div>
  );
}
